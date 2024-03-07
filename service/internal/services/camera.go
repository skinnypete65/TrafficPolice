package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type CameraService interface {
	AddCameraType(cameraType domain.CameraType) error
	RegisterCamera(camera domain.Camera) error
}

type cameraService struct {
	db repository.CameraDB
}

func NewCameraService(db repository.CameraDB) CameraService {
	return &cameraService{db: db}
}

func (s *cameraService) AddCameraType(cameraType domain.CameraType) error {
	id := uuid.New()
	cameraType.ID = id.String()
	return s.db.AddCameraType(cameraType)
}

func (s *cameraService) RegisterCamera(camera domain.Camera) error {
	id := uuid.New()
	camera.ID = id.String()
	return s.db.RegisterCamera(camera)
}
