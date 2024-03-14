package services

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
)

type RatingService interface {
	SetRating(caseDecision domain.CaseDecisionInfo) error
	GetRating() ([]domain.RatingInfo, error)
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

func (s *ratingService) GetRating() ([]domain.RatingInfo, error) {
	rating, err := s.ratingRepo.GetRating()
	if err != nil {
		return nil, err
	}
	if len(rating) == 0 {
		return nil, errs.ErrNoRows
	}

	return rating, nil
}
