package transport

import (
	"fmt"
	"image_storage/internal/services"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

const (
	expertsDir            = "experts"
	expertContentImageKey = "image"
	expertIDPathValue     = "id"
)

type ExpertHandler struct {
	service services.ImgService
}

func NewExpertHandler(service services.ImgService) *ExpertHandler {
	return &ExpertHandler{service: service}
}

func (h *ExpertHandler) UploadExpertImg(w http.ResponseWriter, r *http.Request) {
	expertID := r.PathValue(expertIDPathValue)
	if expertID == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
		return
	}

	file, header, err := parseMultipartForm(r, expertContentImageKey)

	contentType := header.Header.Get(contentTypeKey)
	extension, err := getImgExtension(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgFilePath := fmt.Sprintf("%s/%s.%s", expertsDir, expertID, extension)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error while reading fileBytes: %v\n", fileBytes)
		return
	}

	err = h.service.SaveImg(fileBytes, imgFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Succesfully uploaded image")
}

func (h *ExpertHandler) GetExpertImg(w http.ResponseWriter, r *http.Request) {
	expertID := r.PathValue(expertIDPathValue)
	if expertID == "" {
		http.Error(w, "bad expert id", http.StatusBadRequest)
		return
	}

	pattern := fmt.Sprintf("%s/%s.*", expertsDir, expertID)
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.ServeFile(w, r, files[0])
}
