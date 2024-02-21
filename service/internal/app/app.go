package app

import (
	"TrafficPolice/internal/database/database"
	"TrafficPolice/internal/services/service"
	"TrafficPolice/internal/transport/transport"
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
	cameraService := service.NewCameraService(cameraDB)
	cameraHandler := transport.NewCameraHandler(cameraService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /camera", cameraHandler.AddCameraType)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
