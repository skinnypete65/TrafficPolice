package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
)

type contactInfoRepoPostgres struct {
	conn *pgx.Conn
}

func NewContactInfoRepoPostgres(conn *pgx.Conn) repository.ContactInfoRepo {
	return &contactInfoRepoPostgres{conn: conn}
}

func (r *contactInfoRepoPostgres) InsertContactInfo(m map[string][]*domain.Transport) error {
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

	return r.conn.SendBatch(context.Background(), batch).Close()
}
