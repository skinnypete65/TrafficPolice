package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name CameraService
type CameraService interface {
	AddCameraType(cameraType domain.CameraType) (string, error)
	GetCameraTypeByCameraID(cameraID string) (string, error)
}

type cameraService struct {
	cameraRepo repository.CameraRepo
}

func NewCameraService(db repository.CameraRepo) CameraService {
	return &cameraService{cameraRepo: db}
}

func (s *cameraService) AddCameraType(cameraType domain.CameraType) (string, error) {
	id := uuid.New()
	cameraType.ID = id.String()
	return s.cameraRepo.AddCameraType(cameraType)
}

func (s *cameraService) GetCameraTypeByCameraID(cameraID string) (string, error) {
	return s.cameraRepo.GetCameraTypeByCameraID(cameraID)
}
