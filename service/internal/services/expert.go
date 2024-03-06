package services

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"errors"
)

type ExpertService interface {
	GetExpertByUserID(userID string) (domain.Expert, error)
	GetCaseID(expertID string) (string, error)
}

type expertService struct {
	repo repository.ExpertRepo
}

func NewExpertService(repo repository.ExpertRepo) ExpertService {
	return &expertService{repo: repo}
}

func (s *expertService) GetCaseID(expertID string) (string, error) {
	caseID, err := s.repo.GetLastNotSolvedCase(expertID)
	if err == nil {
		return caseID, nil
	}
	if !errors.Is(err, errs.ErrNoLastNotSolvedCase) {
		return "", err
	}

	return "", err
}

func (s *expertService) GetExpertByUserID(userID string) (domain.Expert, error) {
	return s.repo.GetExpertByUserID(userID)
}
