package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
)

type violationDBPostgres struct {
	conn *pgx.Conn
}

func NewViolationDBPostgres(conn *pgx.Conn) repository.ViolationRepo {
	return &violationDBPostgres{conn: conn}
}

func (db *violationDBPostgres) InsertViolations(violations []*domain.Violation) error {
	batch := &pgx.Batch{}

	query := `INSERT INTO violations (violation_id, violation_name, fine_amount) VALUES ($1, $2, $3)`
	for _, v := range violations {
		batch.Queue(query, v.ID, v.Name, v.FineAmount)
	}

	return db.conn.SendBatch(context.Background(), batch).Close()
}
