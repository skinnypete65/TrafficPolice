package services

import (
	"TrafficPolice/internal/database"
	"TrafficPolice/internal/models"
	"github.com/google/uuid"
)

type CaseService interface {
	AddCase(c *models.Case) error
}

type caseService struct {
	db database.CaseDB
}

func NewCaseService(conn database.CaseDB) CaseService {
	return &caseService{db: conn}
}

func (s *caseService) AddCase(c *models.Case) error {
	id := uuid.New()
	c.ID = id.String()
	c.Transport.ID = "f6e34d03-9c1c-4117-988b-fa61ca6f6c3d"

	return s.db.InsertCase(c)
}
