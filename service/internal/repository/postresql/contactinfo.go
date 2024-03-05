package repository

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
)

type contactInfoDBPostgres struct {
	conn *pgx.Conn
}

func NewContactInfoDBPostgres(conn *pgx.Conn) repository.ContactInfoDB {
	return &contactInfoDBPostgres{conn: conn}
}

func (db *contactInfoDBPostgres) InsertContactInfo(m map[string][]*models.Transport) error {
	batch := &pgx.Batch{}

	personQuery := `INSERT INTO persons (id, phone_num, email, vk_id, tg_id) VALUES ($1, $2, $3, $4, $5)`
	transportQuery := `INSERT INTO transports(transport_id, transport_chars, transport_nums, region, person_id) VALUES ($1, $2, $3, $4, $5)`

	for _, transports := range m {
		if len(transports) > 0 {
			person := transports[0].Person
			batch.Queue(personQuery, person.ID, person.PhoneNum, person.Email, person.VkID, person.TgID)
		}
		for _, t := range transports {
			batch.Queue(transportQuery, t.ID, t.Chars, t.Num, t.Region, t.Person.ID)
		}
	}

	return db.conn.SendBatch(context.Background(), batch).Close()
}
