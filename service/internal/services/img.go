package services

import (
	"log"
	"os"
)

type ImgService interface {
	SaveImg(img []byte, filepath string) error
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
