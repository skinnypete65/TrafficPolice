package app

import (
	"image_storage/internal/services"
	"image_storage/internal/transport"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	imgService := services.NewImgService()
	caseHandler := transport.NewCaseHandler(imgService)

	expertHandler := transport.NewExpertHandler(imgService)

	mux.HandleFunc("POST /case/{id}/img", caseHandler.UploadCaseImg)
	mux.HandleFunc("GET /case/{id}/img", caseHandler.GetCaseImg)

	mux.HandleFunc("POST /expert/{id}/img", expertHandler.UploadExpertImg)
	mux.HandleFunc("GET /expert/{id}/img", expertHandler.GetExpertImg)

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
