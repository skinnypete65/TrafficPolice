package imagereader

import (
	"os"
	"strings"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name ImageReader
type ImageReader interface {
	Read(filePath string) ([]byte, error)
	GetExtension(filePath string) string
}

type imageReader struct {
}

func NewImageReader() ImageReader {
	return &imageReader{}
}

func (r *imageReader) Read(filePath string) ([]byte, error) {
	img, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (r *imageReader) GetExtension(filePath string) string {
	dotIdx := strings.LastIndex(filePath, ".")
	return filePath[dotIdx+1:]
}
