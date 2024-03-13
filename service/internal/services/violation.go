package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type ViolationService interface {
	InsertViolations(violations []*domain.Violation) error
}

type violationService struct {
	db repository.ViolationRepo
}

func NewViolationService(db repository.ViolationRepo) ViolationService {
	return &violationService{db: db}
}

func (s *violationService) InsertViolations(violations []*domain.Violation) error {
	for i := range violations {
		violations[i].ID = uuid.New().String()
	}

	return s.db.InsertViolations(violations)
}
