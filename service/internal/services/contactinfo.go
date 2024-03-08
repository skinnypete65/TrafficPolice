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
	db repository.ContactInfoDB
}

func NewContactInfoService(db repository.ContactInfoDB) ContactInfoService {
	return &contactInfoService{db: db}
}

func (s *contactInfoService) InsertContactInfo(m map[string][]*domain.Transport) error {

	for _, transports := range m {
		personID := uuid.New().String()

		for i := range transports {
			transports[i].Person.ID = personID
			transports[i].ID = uuid.New().String()
		}
	}
	return s.db.InsertContactInfo(m)
}
