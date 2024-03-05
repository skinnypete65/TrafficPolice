package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/models"
)

type CameraDB interface {
	AddCameraType(cameraType models.CameraType) error
	RegisterCamera(camera models.Camera) error
}

type CaseDB interface {
	InsertCase(c *models.Case) error
}

type ContactInfoDB interface {
	InsertContactInfo(m map[string][]*models.Transport) error
}

type ViolationDB interface {
	InsertViolations(violations []*models.Violation) error
}

type AuthRepo interface {
	CheckUserExists(username string) error
	InsertUser(user domain.User) error
	InsertExpert(expert domain.Expert) error
	InsertDirector(director domain.Director) error
	SignIn(username string) (domain.User, error)
}
