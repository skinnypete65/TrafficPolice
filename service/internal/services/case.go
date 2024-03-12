package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type CaseService interface {
	AddCase(c domain.Case) error
}

type caseService struct {
	caseRepo      repository.CaseRepo
	transportRepo repository.TransportRepo
}

func NewCaseService(
	caseRepo repository.CaseRepo,
	transportRepo repository.TransportRepo,
) CaseService {
	return &caseService{
		caseRepo:      caseRepo,
		transportRepo: transportRepo,
	}
}

func (s *caseService) AddCase(c domain.Case) error {
	id := uuid.New()
	c.ID = id.String()
	transportID, err := s.transportRepo.GetTransportID(c.Transport.Chars, c.Transport.Num, c.Transport.Region)

	if err != nil {
		return err
	}
	c.Transport.ID = transportID

	return s.caseRepo.InsertCase(c)
}
