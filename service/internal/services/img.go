package services

import (
	"TrafficPolice/errs"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ImgService interface {
	SaveImg(img []byte, filepath string) error
	GetImgFilePath(dir string, id string) (string, error)
}

type imgServiceLocal struct {
}

func NewImgService() ImgService {
	return &imgServiceLocal{}
}

func (s *imgServiceLocal) SaveImg(img []byte, filepath string) error {
	imgFile, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error while create file to server: %v\n", err)
		return err
	}
	defer imgFile.Close()

	_, err = imgFile.Write(img)
	if err != nil {
		log.Printf("Error while writing tempFile: %v\n", err)
		return err
	}

	return nil
}

func (s *imgServiceLocal) GetImgFilePath(dir string, id string) (string, error) {
	pattern := fmt.Sprintf("%s/%s.*", dir, id)
	files, err := filepath.Glob(pattern)
	if len(files) == 0 {
		return "", errs.ErrNoImage
	}
	if err != nil {
		return "", err
	}
	return files[0], nil
}
