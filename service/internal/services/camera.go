package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type CameraService interface {
	AddCameraType(cameraType domain.CameraType) error
	GetCameraTypeByCameraID(cameraID string) (string, error)
}

type cameraService struct {
	cameraRepo repository.CameraRepo
}

func NewCameraService(db repository.CameraRepo) CameraService {
	return &cameraService{cameraRepo: db}
}

func (s *cameraService) AddCameraType(cameraType domain.CameraType) error {
	id := uuid.New()
	cameraType.ID = id.String()
	return s.cameraRepo.AddCameraType(cameraType)
}

func (s *cameraService) GetCameraTypeByCameraID(cameraID string) (string, error) {
	return s.cameraRepo.GetCameraTypeByCameraID(cameraID)
}
