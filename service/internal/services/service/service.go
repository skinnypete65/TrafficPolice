package service

import (
	"TrafficPolice/internal/database"
	"TrafficPolice/internal/models"
	"github.com/google/uuid"
	"log"
)

type CameraService interface {
	AddCameraType(cameraType models.CameraType) error
	RegisterCamera(camera models.Camera) error
}

type cameraService struct {
	db database.CameraDB
}

func NewCameraService(db database.CameraDB) CameraService {
	return &cameraService{db: db}
}

func (s *cameraService) AddCameraType(cameraType models.CameraType) error {
	id := uuid.New()
	cameraType.ID = id.String()

	log.Println(cameraType)
	log.Println(id)

	return s.db.AddCameraType(cameraType)
}

func (s *cameraService) RegisterCamera(camera models.Camera) error {
	return nil
}
