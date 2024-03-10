package app

import (
	"TrafficPolice/internal/config"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository/postresql"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport/rest"
	"TrafficPolice/internal/transport/rest/middlewares"
	"TrafficPolice/internal/validation"
	"TrafficPolice/pkg/rabbitmq"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
)

func Run() {
	cfg, err := config.ParseConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(context.Background())

	rabbitMQConn, err := rabbitmq.NewRabbitMQConn()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQConn.Close()

	validate := newValidate()

	tokenManager, _ := tokens.NewTokenManager(cfg.SigningKey)
	imgService := services.NewImgService()

	paginationRepo := repository.NewPaginationRepoPostgres(dbConn)
	paginationService := services.NewPaginationService(paginationRepo)

	caseRepo := repository.NewCaseRepoPostgres(dbConn)
	caseService := services.NewCaseService(caseRepo)
	caseHandler := rest.NewCaseHandler(caseService, imgService)

	cameraDB := repository.NewCameraRepoPostgres(dbConn)
	cameraService := services.NewCameraService(cameraDB)
	cameraHandler := rest.NewCameraHandler(cameraService, validate)

	contactInfoDB := repository.NewContactInfoDBPostgres(dbConn)
	contactService := services.NewContactInfoService(contactInfoDB)
	contactInfoHandler := rest.NewContactInfoHandler(contactService)

	violationDB := repository.NewViolationDBPostgres(dbConn)
	violationService := services.NewViolationService(violationDB)
	violationHandler := rest.NewViolationHandler(violationService)

	authRepo := repository.NewAuthRepoPostgres(dbConn)
	authService := services.NewAuthService(authRepo, tokenManager, cfg.PassSalt)
	authHandler := rest.NewAuthHandler(authService, validate)

	expertRepo := repository.NewExpertRepoPostgres(dbConn)
	expertService := services.NewExpertService(expertRepo, caseRepo, cfg.Consensus)
	expertHandler := rest.NewExpertHandler(imgService, expertService)

	authMiddleware := middlewares.NewAuthMiddleware(tokenManager, expertService)

	trainingRepo := repository.NewTrainingRepoPostgres(dbConn)
	trainingService := services.NewTrainingService(trainingRepo)
	trainingHandler := rest.NewTrainingHandler(trainingService, paginationService, validate)

	mux := http.NewServeMux()

	registerDirectors(cfg, authService)

	mux.Handle("POST /camera/type",
		authMiddleware.IdentifyRole(http.HandlerFunc(cameraHandler.AddCameraType), domain.DirectorRole),
	)
	mux.Handle("POST /camera",
		authMiddleware.IdentifyRole(http.HandlerFunc(cameraHandler.RegisterCamera), domain.DirectorRole),
	)

	mux.HandleFunc("POST /case", caseHandler.AddCase)
	mux.Handle("POST /case/{id}/img", http.HandlerFunc(caseHandler.UploadCaseImg))
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

	port := fmt.Sprintf(":%d", cfg.ServerPort)
	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

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
