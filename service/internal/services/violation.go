package services

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type ViolationService interface {
	InsertViolations(violations []*models.Violation) error
}

type violationService struct {
	db repository.ViolationDB
}

func NewViolationService(db repository.ViolationDB) ViolationService {
	return &violationService{db: db}
}

func (s *violationService) InsertViolations(violations []*models.Violation) error {
	for i := range violations {
		violations[i].ID = uuid.New().String()
	}

	return s.db.InsertViolations(violations)
}
