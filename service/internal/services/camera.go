package services

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type CameraService interface {
	AddCameraType(cameraType models.CameraType) error
	RegisterCamera(camera models.Camera) error
}

type cameraService struct {
	db repository.CameraDB
}

func NewCameraService(db repository.CameraDB) CameraService {
	return &cameraService{db: db}
}

func (s *cameraService) AddCameraType(cameraType models.CameraType) error {
	id := uuid.New()
	cameraType.ID = id.String()
	return s.db.AddCameraType(cameraType)
}

func (s *cameraService) RegisterCamera(camera models.Camera) error {
	id := uuid.New()
	camera.ID = id.String()
	return s.db.RegisterCamera(camera)
}
