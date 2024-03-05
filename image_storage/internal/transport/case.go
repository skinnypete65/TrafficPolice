package transport

import (
	"fmt"
	"image_storage/internal/services"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var allowedExtensions = map[string]struct{}{
	"png":  {},
	"jpeg": {},
	"jpg":  {},
}

const (
	contentTypeKey = "Content-Type"
	contentImage   = "image"
	casesDir       = "cases"
)

type CaseHandler struct {
	service services.ImgService
}

func NewCaseHandler(service services.ImgService) *CaseHandler {
	return &CaseHandler{service: service}
}

func (h *CaseHandler) UploadCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
		return
	}

	file, header, err := parseMultipartForm(r, "image")

	contentType := header.Header.Get(contentTypeKey)
	extension, err := getImgExtension(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgFilePath := fmt.Sprintf("%s/%s.%s", casesDir, caseID, extension)
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

func (h *CaseHandler) GetCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		http.Error(w, "bad case id", http.StatusBadRequest)
		return
	}

	pattern := fmt.Sprintf("%s/%s.*", casesDir, caseID)
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.ServeFile(w, r, files[0])
}
