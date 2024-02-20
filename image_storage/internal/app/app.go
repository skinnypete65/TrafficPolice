package app

import (
	"image_storage/internal/transport"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /case/{id}/img", transport.UploadCaseImg)

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
