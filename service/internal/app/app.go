package app

import (
	"TrafficPolice/internal/repository/postresql"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport"
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

	cameraDB := repository.NewCameraRepoPostgres(conn)
	cameraService := services.NewCameraService(cameraDB)
	cameraHandler := transport.NewCameraHandler(cameraService)

	caseDB := repository.NewCaseDBPostgres(conn)
	caseService := services.NewCaseService(caseDB)
	caseHandler := transport.NewCaseHandler(caseService)

	contactInfoDB := repository.NewContactInfoDBPostgres(conn)
	contactService := services.NewContactInfoService(contactInfoDB)
	contactInfoHandler := transport.NewContactInfoHandler(contactService)

	violationDB := repository.NewViolationDBPostgres(conn)
	violationService := services.NewViolationService(violationDB)
	violationHandler := transport.NewViolationHandler(violationService)

	authRepo := repository.NewAuthRepoPostgres(conn)
	authService := services.NewAuthService(authRepo)
	authHandler := transport.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /camera/type", cameraHandler.AddCameraType)
	mux.HandleFunc("POST /camera", cameraHandler.RegisterCamera)

	mux.HandleFunc("POST /case", caseHandler.AddCase)

	mux.HandleFunc("POST /contact_info", contactInfoHandler.InsertContactInfo)

	mux.HandleFunc("POST /violations", violationHandler.InsertViolations)

	mux.HandleFunc("POST /auth/sign_up", authHandler.SignUp)
	mux.HandleFunc("POST /auth/sign_in", authHandler.SignIn)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
