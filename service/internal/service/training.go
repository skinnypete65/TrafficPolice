package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
)

type TrainingService interface {
	GetSolvedCasesByParams(params domain.SolvedCasesParams, paginationParams domain.PaginationParams) ([]domain.Case, error)
}

type trainingService struct {
	trainingRepo repository.TrainingRepo
}

func NewTrainingService(trainingRepo repository.TrainingRepo) TrainingService {
	return &trainingService{
		trainingRepo: trainingRepo,
	}
}

func (s *trainingService) GetSolvedCasesByParams(
	params domain.SolvedCasesParams,
	paginationParams domain.PaginationParams,
) ([]domain.Case, error) {
	return s.trainingRepo.GetSolvedCasesByParams(params, paginationParams)
}
