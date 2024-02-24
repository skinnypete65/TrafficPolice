package services

import (
	"TrafficPolice/internal/database"
	"TrafficPolice/internal/models"
	"github.com/google/uuid"
)

type ContactInfoService interface {
	InsertContactInfo(m map[string][]*models.Transport) error
}

type contactInfoService struct {
	db database.ContactInfoDB
}

func NewContactInfoService(db database.ContactInfoDB) ContactInfoService {
	return &contactInfoService{db: db}
}

func (s *contactInfoService) InsertContactInfo(m map[string][]*models.Transport) error {

	for _, transports := range m {
		personID := uuid.New().String()

		for i := range transports {
			transports[i].Person.ID = personID
			transports[i].ID = uuid.New().String()
		}
	}
	return s.db.InsertContactInfo(m)
}
