package repository

import (
	"TrafficPolice/internal/domain"
	"github.com/google/uuid"
	"time"
)

type CameraRepo interface {
	AddCameraType(cameraType domain.CameraType) (string, error)
	GetCameraTypeByCameraID(cameraID string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name CaseRepo
type CaseRepo interface {
	InsertCase(c domain.Case) (string, error)
	GetCaseByID(caseID string) (domain.Case, error)
	GetCaseWithPersonInfo(caseID string) (domain.Case, error)
	SetCaseFineDecision(caseID string, fineDecision bool, solvedAt time.Time) error
	UpdateCaseRequiredSkill(caseID string, requiredSkill int) error
}

type ContactInfoRepo interface {
	InsertContactInfo(m map[string][]*domain.Transport) error
}

type ViolationRepo interface {
	InsertViolations(violations []*domain.Violation) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name AuthRepo
type AuthRepo interface {
	CheckUserExists(username string) bool
	InsertUser(user domain.UserInfo) error
	InsertExpert(expert domain.Expert) error
	InsertCamera(camera domain.Camera, userID uuid.UUID) (string, error)
	InsertDirector(director domain.Director) error
	SignIn(username string) (domain.UserInfo, error)
	ConfirmExpert(data domain.ConfirmExpert) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name ExpertRepo
type ExpertRepo interface {
	GetLastNotSolvedCaseID(expertID string) (string, error)
	GetExpertByUserID(userID string) (domain.Expert, error)
	GetNotSolvedCase(expert domain.Expert) (domain.Case, error)
	InsertNotSolvedCase(solvedCase domain.ExpertCase) error
	SetCaseDecision(decision domain.Decision) error
	GetCaseFineDecisions(caseID string, competenceSkill int) (domain.FineDecisions, error)
	GetExpertsCountBySkill(competenceSkill int) (int, error)
}

type TrainingRepo interface {
	GetSolvedCasesByParams(params domain.SolvedCasesParams, paginationParams domain.PaginationParams) ([]domain.Case, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name PaginationRepo
type PaginationRepo interface {
	GetRecordsCount(table string) (int, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name TransportRepo
type TransportRepo interface {
	GetTransportID(chars string, num string, region string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name RatingRepo
type RatingRepo interface {
	GetSolvedCaseDecisions(caseDecision domain.CaseDecisionInfo) ([]domain.ExpertCaseDecision, error)
	SetRating(decisions []domain.ExpertCaseDecision) error
	InsertExpertId(expertID string) error
	GetRating() ([]domain.RatingInfo, error)
	GetExpertsRating(minSolvedCases int) ([]domain.ExpertRating, error)
	UpdateCompetenceSkills(infos []domain.UpdateCompetenceSkill) error
	ClearRating() error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name DirectorRepo
type DirectorRepo interface {
	GetCase(caseID string) (domain.CaseStatus, error)
	GetExpertIntervalCases(
		expertID string,
		startDate time.Time,
		endDate time.Time) (map[domain.Date][]domain.IntervalCase, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name CheckerRepo
type CheckerRepo interface {
	CheckExpertExists(expertID string) (bool, error)
}
