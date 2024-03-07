package repository

import (
	"TrafficPolice/internal/domain"
)

type CameraDB interface {
	AddCameraType(cameraType domain.CameraType) error
	RegisterCamera(camera domain.Camera) error
}

type CaseRepo interface {
	InsertCase(c *domain.Case) error
	GetCaseByID(caseID string) (domain.Case, error)
}

type ContactInfoDB interface {
	InsertContactInfo(m map[string][]*domain.Transport) error
}

type ViolationDB interface {
	InsertViolations(violations []*domain.Violation) error
}

type AuthRepo interface {
	CheckUserExists(username string) error
	InsertUser(user domain.UserInfo) error
	InsertExpert(expert domain.Expert) error
	InsertDirector(director domain.Director) error
	SignIn(username string) (domain.UserInfo, error)
	ConfirmExpert(data domain.ConfirmExpert) error
}

type ExpertRepo interface {
	GetLastNotSolvedCaseID(expertID string) (string, error)
	GetExpertByUserID(userID string) (domain.Expert, error)
	GetNotSolvedCase(expert domain.Expert) (domain.Case, error)
	InsertNotSolvedCase(solvedCase domain.SolvedCase) error
	SetCaseDecision(decision domain.Decision) error
}
