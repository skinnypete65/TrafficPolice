package app

import (
	_ "TrafficPolice/docs"
	"TrafficPolice/internal/config"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport/rabbitmq"
	"TrafficPolice/internal/transport/rest/middlewares"
	"TrafficPolice/internal/validation"
	"TrafficPolice/pkg/imagereader"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

const (
	serviceConfigPath = "service_config.yaml"
	defaultMinExperts = 3
	serverPort        = ":8080"
)

func Run() {
	cfg, err := config.ParseConfig(serviceConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Rating.MinExperts < defaultMinExperts {
		log.Fatalf("config min experts must be greater or equal %d, but got: %d",
			defaultMinExperts, cfg.Rating.MinExperts,
		)
	}

	// Init database
	dbConnString := setupDBConnString(cfg)
	dbConn, err := pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(context.Background())

	runMigrations(dbConnString)

	// Init RabbitMQ
	mQConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer mQConn.Close()

	// Init handlers, service, repos (Clean architecture)
	tokenManager, err := tokens.NewTokenManager(cfg.SigningKey)
	if err != nil {
		log.Fatal(err)
	}
	finePublisher := setupFinePublisher(mQConn)
	validate := newValidate()
	imageReader := imagereader.NewImageReader()

	converters := newConverters()
	repos := newRepos(dbConn)
	services := newServices(repos, tokenManager, cfg)
	handlers := newHandlers(services, converters, validate, finePublisher, imageReader)

	authMiddleware := middlewares.NewAuthMiddleware(tokenManager, services.expert)
	serveMuxInit := newServeMuxInit(handlers, authMiddleware)
	mux := serveMuxInit.Init()

	// Run logic
	registerDirectors(cfg, services.auth)

	done := make(chan struct{})
	go services.rating.RunReportPeriod(done)

	// Run Server
	server := http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	log.Printf("Run server on %s\n", serverPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	done <- struct{}{}

}

func runMigrations(dbUrl string) {
	log.Printf("Run migrations on %s\n", dbUrl)
	m, err := migrate.New("file://migrations", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate no change")
	} else if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrate ran successfully")
}

func registerDirectors(cfg *config.Config, authService service.AuthService) {
	users := make([]domain.UserInfo, len(cfg.Directors))

	for i, d := range cfg.Directors {
		users[i] = domain.UserInfo{Username: d.Username, Password: d.Password}
	}

	err := authService.RegisterDirectors(users)
	if err != nil {
		log.Fatal(err)
	}
}

func newValidate() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("is_date_only", validation.IsDateOnly)
	if err != nil {
		log.Println(err)
	}

	return validate
}

func setupDBConnString(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}

func setupFinePublisher(mqConn *amqp.Connection) *rabbitmq.FinePublisherRabbitMQ {
	finePublisher, err := rabbitmq.NewFinePublisher(mqConn)
	if err != nil {
		log.Fatal(err)
	}
	err = finePublisher.SetupExchangeAndQueue(
		rabbitmq.ExchangeParams{
			Name:       rabbitmq.FineExchange,
			Kind:       rabbitmq.Fanout,
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		}, rabbitmq.QueueParams{
			Name:       rabbitmq.FineQueue,
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		rabbitmq.BindingParams{
			Queue:    rabbitmq.FineQueue,
			Key:      "",
			Exchange: rabbitmq.FineExchange,
			NoWait:   false,
			Args:     nil,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return finePublisher
}
