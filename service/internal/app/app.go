package app

import (
	"TrafficPolice/internal/camera"
	"TrafficPolice/internal/config"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository/postresql"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport/rabbitmq"
	"TrafficPolice/internal/transport/rest"
	"TrafficPolice/internal/transport/rest/middlewares"
	"TrafficPolice/internal/validation"
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
)

func Run() {
	cfg, err := config.ParseConfig(serviceConfigPath)
	if err != nil {
		log.Fatal(err)
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

	finePublisher := setupFinePublisher(mQConn)

	validate := newValidate()

	// Init handlers, services, repos (Clean architecture)
	tokenManager, _ := tokens.NewTokenManager(cfg.SigningKey)
	imgService := services.NewImgService()

	authRepo := repository.NewAuthRepoPostgres(dbConn)
	authService := services.NewAuthService(authRepo, tokenManager, cfg.PassSalt)
	authHandler := rest.NewAuthHandler(authService, validate)

	paginationRepo := repository.NewPaginationRepoPostgres(dbConn)
	paginationService := services.NewPaginationService(paginationRepo)

	cameraRepo := repository.NewCameraRepoPostgres(dbConn)
	cameraService := services.NewCameraService(cameraRepo)
	cameraHandler := rest.NewCameraHandler(cameraService, authService, validate)

	transportRepo := repository.NewTransportRepoPostgres(dbConn)
	caseRepo := repository.NewCaseRepoPostgres(dbConn)
	caseService := services.NewCaseService(caseRepo, transportRepo)
	caseConverter := converter.NewCaseConverter()
	cameraParser := camera.NewParser(cameraService)
	caseHandler := rest.NewCaseHandler(caseService, imgService, cameraService, caseConverter, cameraParser)

	contactInfoDB := repository.NewContactInfoRepoPostgres(dbConn)
	contactService := services.NewContactInfoService(contactInfoDB)
	contactInfoHandler := rest.NewContactInfoHandler(contactService)

	violationDB := repository.NewViolationDBPostgres(dbConn)
	violationService := services.NewViolationService(violationDB)
	violationHandler := rest.NewViolationHandler(violationService)

	expertRepo := repository.NewExpertRepoPostgres(dbConn)
	expertService := services.NewExpertService(expertRepo, caseRepo, cfg.Consensus)
	expertHandler := rest.NewExpertHandler(imgService, expertService, finePublisher, caseConverter)

	authMiddleware := middlewares.NewAuthMiddleware(tokenManager, expertService)

	trainingRepo := repository.NewTrainingRepoPostgres(dbConn)
	trainingService := services.NewTrainingService(trainingRepo)
	paginationConverter := converter.NewPaginationConverter()
	trainingHandler := rest.NewTrainingHandler(trainingService, paginationService, validate,
		caseConverter, paginationConverter)

	registerDirectors(cfg, authService)

	// Setup Routes
	mux := http.NewServeMux()

	mux.Handle("POST /camera/type",
		authMiddleware.IdentifyRole(http.HandlerFunc(cameraHandler.AddCameraType), domain.DirectorRole),
	)
	mux.Handle("POST /camera",
		authMiddleware.IdentifyRole(http.HandlerFunc(cameraHandler.RegisterCamera), domain.DirectorRole),
	)

	mux.Handle("POST /case",
		authMiddleware.IdentifyRole(http.HandlerFunc(caseHandler.AddCase), domain.CameraRole),
	)
	mux.Handle("POST /case/{id}/img",
		authMiddleware.IdentifyRole(http.HandlerFunc(caseHandler.UploadCaseImg), domain.CameraRole),
	)
	mux.Handle("GET /case/{id}/img",
		authMiddleware.IdentifyRole(
			authMiddleware.IsExpertConfirmed(http.HandlerFunc(caseHandler.GetCaseImg)),
			domain.DirectorRole, domain.ExpertRole,
		),
	)

	mux.Handle("POST /contact_info",
		authMiddleware.IdentifyRole(http.HandlerFunc(contactInfoHandler.InsertContactInfo), domain.DirectorRole),
	)

	mux.Handle("POST /violations",
		authMiddleware.IdentifyRole(http.HandlerFunc(violationHandler.InsertViolations), domain.DirectorRole),
	)

	mux.HandleFunc("POST /auth/sign_up", authHandler.SignUp)
	mux.HandleFunc("POST /auth/sign_in", authHandler.SignIn)
	mux.Handle("POST /auth/confirm/expert",
		authMiddleware.IdentifyRole(http.HandlerFunc(authHandler.ConfirmExpert), domain.DirectorRole),
	)

	mux.Handle("POST /expert/{id}/img", authMiddleware.IdentifyRole(
		authMiddleware.IsExpertConfirmed(http.HandlerFunc(expertHandler.UploadExpertImg)),
		domain.DirectorRole, domain.ExpertRole),
	)
	mux.Handle("GET /expert/{id}/img",
		authMiddleware.IdentifyRole(
			authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(expertHandler.GetExpertImg),
			),
			domain.DirectorRole, domain.ExpertRole,
		),
	)

	mux.Handle("GET /expert/get_case",
		authMiddleware.IdentifyRole(
			authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(expertHandler.GetCaseForExpert),
			),
			domain.ExpertRole,
		),
	)
	mux.Handle("POST /expert/decision",
		authMiddleware.IdentifyRole(
			authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(expertHandler.SetCaseDecision),
			),
			domain.ExpertRole,
		),
	)

	mux.Handle("POST /expert/training",
		authMiddleware.IdentifyRole(
			authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(trainingHandler.GetSolvedCasesByParams),
			),
			domain.ExpertRole,
		),
	)

	// Run Server
	port := fmt.Sprintf(":%d", cfg.ServerPort)
	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Run server on %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

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

func registerDirectors(cfg *config.Config, authService services.AuthService) {
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

func setupFinePublisher(mqConn *amqp.Connection) *rabbitmq.FinePublisher {
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
