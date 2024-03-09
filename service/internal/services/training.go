package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
)

type TrainingService interface {
	GetSolvedCasesByParams(params domain.SolvedCasesParams) ([]domain.Case, error)
}

type trainingService struct {
	trainingRepo repository.TrainingRepo
}

func NewTrainingService(trainingRepo repository.TrainingRepo) TrainingService {
	return &trainingService{
		trainingRepo: trainingRepo,
	}
}

func (s *trainingService) GetSolvedCasesByParams(params domain.SolvedCasesParams) ([]domain.Case, error) {
	return s.trainingRepo.GetSolvedCasesByParams(params)
}
