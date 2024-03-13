package services

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"github.com/google/uuid"
)

type ContactInfoService interface {
	InsertContactInfo(m map[string][]*domain.Transport) error
}

type contactInfoService struct {
	repo repository.ContactInfoRepo
}

func NewContactInfoService(repo repository.ContactInfoRepo) ContactInfoService {
	return &contactInfoService{repo: repo}
}

func (s *contactInfoService) InsertContactInfo(m map[string][]*domain.Transport) error {
	for _, transports := range m {
		personID := uuid.New().String()

		for i := range transports {
			transports[i].Person.ID = personID
			transports[i].ID = uuid.New().String()
		}
	}
	return s.repo.InsertContactInfo(m)
}
