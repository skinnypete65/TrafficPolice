package repository

import (
	"TrafficPolice/internal/domain"
	"github.com/google/uuid"
)

type CameraRepo interface {
	AddCameraType(cameraType domain.CameraType) (string, error)
	GetCameraTypeByCameraID(cameraID string) (string, error)
}

type CaseRepo interface {
	InsertCase(c domain.Case) (string, error)
	GetCaseByID(caseID string) (domain.Case, error)
	GetCaseWithPersonInfo(caseID string) (domain.Case, error)
	SetCaseFineDecision(caseID string, fineDecision bool) error
	UpdateCaseRequiredSkill(caseID string, requiredSkill int) error
}

type ContactInfoRepo interface {
	InsertContactInfo(m map[string][]*domain.Transport) error
}

type ViolationRepo interface {
	InsertViolations(violations []*domain.Violation) error
}

type AuthRepo interface {
	CheckUserExists(username string) bool
	InsertUser(user domain.UserInfo) error
	InsertExpert(expert domain.Expert) error
	InsertCamera(camera domain.Camera, userID uuid.UUID) (string, error)
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
	GetCaseFineDecisions(caseID string) (domain.FineDecisions, error)
	GetExpertsCountBySkill(competenceSkill int) (int, error)
}

type DirectorRepo interface {
	InsertDirectors(directors []domain.Director) error
}

type TrainingRepo interface {
	GetSolvedCasesByParams(params domain.SolvedCasesParams, paginationParams domain.PaginationParams) ([]domain.Case, error)
}

type PaginationRepo interface {
	GetRecordsCount(table string) (int, error)
}

type TransportRepo interface {
	GetTransportID(chars string, num string, region string) (string, error)
}

type RatingRepo interface {
	GetSolvedCaseDecisions(caseDecision domain.CaseDecisionInfo) ([]domain.SolvedCaseDecision, error)
	SetRating(decisions []domain.SolvedCaseDecision) error
	InsertExpertId(expertID string) error
}
