package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
)

type RatingService interface {
	SetRating(caseDecision domain.CaseDecisionInfo) error
}

type ratingService struct {
	ratingRepo repository.RatingRepo
}

func NewRatingService(ratingRepo repository.RatingRepo) RatingService {
	return &ratingService{
		ratingRepo: ratingRepo,
	}
}

func (s *ratingService) SetRating(caseDecision domain.CaseDecisionInfo) error {
	solvedDecisions, err := s.ratingRepo.GetSolvedCaseDecisions(caseDecision)
	if err != nil {
		return err
	}

	return s.ratingRepo.SetRating(solvedDecisions)
}
