package transport

import (
	"fmt"
	"image_storage/internal/services"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]struct{}{
	"png":  {},
	"jpeg": {},
	"jpg":  {},
}

const (
	contentTypeKey = "Content-Type"
	contentImage   = "image"
)

type ImgHandler struct {
	service services.ImgService
}

func NewImgHandler(service services.ImgService) *ImgHandler {
	return &ImgHandler{service: service}
}

func (h *ImgHandler) UploadCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
		return
	}

	// parse input, type multipart/form-data
	// 10 MB
	maxMemory := int64(10 << 20)

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		log.Printf("Error while ParseMultipartForm: %v", err)
		return
	}

	// retrieve file from posted form-data
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error retrieving file from form-data: %v\n", err)
		return
	}
	defer file.Close()

	log.Printf("Uploaded file: %+v\n", handler.Filename)
	log.Printf("File size: %+v\n", handler.Size)
	log.Printf("MIME header: %+v\n", handler.Header)

	contentType := handler.Header.Get(contentTypeKey)
	extension, err := getImgExtension(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write file to server
	imgFilePath := fmt.Sprintf("images/%s.%s", caseID, extension)
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

func getImgExtension(contentType string) (string, error) {
	slashIdx := strings.Index(contentType, "/")
	if slashIdx == -1 {
		return "", fmt.Errorf("bad content type")
	}

	fileType := contentType[:slashIdx]
	extension := contentType[slashIdx+1:]

	if fileType != contentImage {
		return "", fmt.Errorf("file type %s is not image", fileType)
	}

	_, isAllowed := allowedExtensions[extension]
	if isAllowed {
		return extension, nil
	} else {
		return "", fmt.Errorf("extension %s is not allowed", extension)
	}
}

func (h *ImgHandler) GetCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("id")
	if caseID == "" {
		http.Error(w, "bad case id", http.StatusBadRequest)
		return
	}

	pattern := fmt.Sprintf("images/%s.*", caseID)
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.ServeFile(w, r, files[0])
}
