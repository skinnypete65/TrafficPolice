package services

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"errors"
	"github.com/google/uuid"
)

type ExpertService interface {
	GetExpertByUserID(userID string) (domain.Expert, error)
	GetCase(userID string) (domain.Case, error)
}

type expertService struct {
	expertRepo repository.ExpertRepo
	caseRepo   repository.CaseRepo
}

func NewExpertService(expertRepo repository.ExpertRepo, caseRepo repository.CaseRepo) ExpertService {
	return &expertService{
		expertRepo: expertRepo,
		caseRepo:   caseRepo,
	}
}

func (s *expertService) GetCase(userID string) (domain.Case, error) {
	expert, err := s.expertRepo.GetExpertByUserID(userID)
	if err != nil {
		return domain.Case{}, err
	}

	caseID, err := s.expertRepo.GetLastNotSolvedCaseID(expert.ID)
	if err == nil {
		return s.caseRepo.GetCaseByID(caseID)
	}
	if !errors.Is(err, errs.ErrNoLastNotSolvedCase) {
		return domain.Case{}, err
	}

	notSolvedCase, err := s.expertRepo.GetNotSolvedCase(expert)
	if err != nil {
		return domain.Case{}, err
	}

	err = s.expertRepo.InsertNotSolvedCase(domain.SolvedCase{
		SolvedCaseID:  uuid.New().String(),
		ExpertID:      expert.ID,
		CaseID:        notSolvedCase.ID,
		IsExpertSolve: false,
		FineDecision:  false,
	})
	if err != nil {
		return domain.Case{}, err
	}

	return notSolvedCase, err
}

func (s *expertService) GetExpertByUserID(userID string) (domain.Expert, error) {
	return s.expertRepo.GetExpertByUserID(userID)
}
