package repository

import (
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type transportRepoPostgres struct {
	conn *pgx.Conn
}

func NewTransportRepoPostgres(conn *pgx.Conn) repository.TransportRepo {
	return &transportRepoPostgres{
		conn: conn,
	}
}

const getTransportIDQuery = `SELECT transport_id FROM transports
WHERE transport_chars = $1 and transport_nums = $2 and region = $3`

func (r *transportRepoPostgres) GetTransportID(chars string, num string, region string) (string, error) {
	row := r.conn.QueryRow(context.Background(), getTransportIDQuery, chars, num, region)

	var transportID string
	err := row.Scan(&transportID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.ErrNoTransport
		}
	}
	return transportID, nil
}
