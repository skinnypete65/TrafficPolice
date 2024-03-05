package transport

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
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

func parseMultipartForm(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	// parse input, type multipart/form-data
	// 10 MB
	maxMemory := int64(10 << 20)

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		log.Printf("Error while ParseMultipartForm: %v", err)
		return nil, nil, err
	}

	// retrieve file from posted form-data
	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving file from form-data: %v\n", err)
	}

	log.Printf("Uploaded file: %+v\n", header.Filename)
	log.Printf("File size: %+v\n", header.Size)
	log.Printf("MIME header: %+v\n", header.Header)

	return file, header, nil
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
