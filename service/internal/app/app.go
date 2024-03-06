package app

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository/postresql"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport"
	"TrafficPolice/internal/transport/middlewares"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
)

func Run() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	tokenManager, _ := tokens.NewTokenManager("sign")
	imgService := services.NewImgService()

	caseDB := repository.NewCaseDBPostgres(conn)
	caseService := services.NewCaseService(caseDB)
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
	authService := services.NewAuthService(authRepo, tokenManager)
	authHandler := transport.NewAuthHandler(authService)

	expertRepo := repository.NewExpertRepoPostgres(conn)
	expertService := services.NewExpertService(expertRepo)
	expertHandler := transport.NewExpertHandler(imgService, expertService)

	authMiddleware := middlewares.NewAuthMiddleware(tokenManager, expertService)

	mux := http.NewServeMux()

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
			authMiddleware.IsExpertConfirmed(http.HandlerFunc(expertHandler.GetExpertImg)),
			domain.DirectorRole, domain.ExpertRole,
		),
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
