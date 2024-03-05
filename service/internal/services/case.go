package services

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type CaseService interface {
	AddCase(c *models.Case) error
}

type caseService struct {
	db repository.CaseDB
}

func NewCaseService(conn repository.CaseDB) CaseService {
	return &caseService{db: conn}
}

func (s *caseService) AddCase(c *models.Case) error {
	id := uuid.New()
	c.ID = id.String()
	c.Transport.ID = "f6e34d03-9c1c-4117-988b-fa61ca6f6c3d"

	return s.db.InsertCase(c)
}
