package app

import (
	"TrafficPolice/internal/config"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository/postresql"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport"
	"TrafficPolice/internal/transport/middlewares"
	"context"
	"fmt"
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

	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	tokenManager, _ := tokens.NewTokenManager(cfg.SigningKey)
	imgService := services.NewImgService()

	caseRepo := repository.NewCaseRepoPostgres(conn)
	caseService := services.NewCaseService(caseRepo)
	caseHandler := transport.NewCaseHandler(caseService, imgService)

	cameraDB := repository.NewCameraRepoPostgres(conn)
	cameraService := services.NewCameraService(cameraDB)
	cameraHandler := transport.NewCameraHandler(cameraService)

	contactInfoDB := repository.NewContactInfoDBPostgres(conn)
	contactService := services.NewContactInfoService(contactInfoDB)
	contactInfoHandler := transport.NewContactInfoHandler(contactService)

	violationDB := repository.NewViolationDBPostgres(conn)
	violationService := services.NewViolationService(violationDB)
	violationHandler := transport.NewViolationHandler(violationService)

	authRepo := repository.NewAuthRepoPostgres(conn)
	authService := services.NewAuthService(authRepo, tokenManager, cfg.PassSalt)
	authHandler := transport.NewAuthHandler(authService)

	expertRepo := repository.NewExpertRepoPostgres(conn)
	expertService := services.NewExpertService(expertRepo, caseRepo, cfg.Consensus)
	expertHandler := transport.NewExpertHandler(imgService, expertService)

	authMiddleware := middlewares.NewAuthMiddleware(tokenManager, expertService)

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
