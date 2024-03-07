package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type CaseService interface {
	AddCase(c *domain.Case) error
}

type caseService struct {
	repo repository.CaseRepo
}

func NewCaseService(conn repository.CaseRepo) CaseService {
	return &caseService{repo: conn}
}

func (s *caseService) AddCase(c *domain.Case) error {
	id := uuid.New()
	c.ID = id.String()
	c.Transport.ID = "f6e34d03-9c1c-4117-988b-fa61ca6f6c3d"

	return s.repo.InsertCase(c)
}
