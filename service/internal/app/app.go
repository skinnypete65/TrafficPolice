package app

import (
	"TrafficPolice/internal/database/postresql"
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

	cameraDB := database.NewCameraDBPostgres(conn)
	cameraService := services.NewCameraService(cameraDB)
	cameraHandler := transport.NewCameraHandler(cameraService)

	caseDB := database.NewCaseDBPostgres(conn)
	caseService := services.NewCaseService(caseDB)
	caseHandler := transport.NewCaseHandler(caseService)

	contactInfoDB := database.NewContactInfoDBPostgres(conn)
	contactService := services.NewContactInfoService(contactInfoDB)
	contactInfoHandler := transport.NewContactInfoHandler(contactService)

	violationDB := database.NewViolationDBPostgres(conn)
	violationService := services.NewViolationService(violationDB)
	violationHandler := transport.NewViolationHandler(violationService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /camera/type", cameraHandler.AddCameraType)
	mux.HandleFunc("POST /camera", cameraHandler.RegisterCamera)

	mux.HandleFunc("POST /case", caseHandler.AddCase)

	mux.HandleFunc("POST /contact_info", contactInfoHandler.InsertContactInfo)

	mux.HandleFunc("POST /violations", violationHandler.InsertViolations)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
