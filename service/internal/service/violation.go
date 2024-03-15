package service

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type ViolationService interface {
	InsertViolations(violations []*domain.Violation) error
}

type violationService struct {
	repo repository.ViolationRepo
}

func NewViolationService(db repository.ViolationRepo) ViolationService {
	return &violationService{repo: db}
}

func (s *violationService) InsertViolations(violations []*domain.Violation) error {
	for i := range violations {
		violations[i].ID = uuid.New().String()
	}

	return s.repo.InsertViolations(violations)
}
