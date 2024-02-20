package app

import (
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
