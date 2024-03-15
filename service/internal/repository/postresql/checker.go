package repository

import (
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type checkerRepoPostgres struct {
	conn *pgx.Conn
}

func NewCheckerRepoPostgres(conn *pgx.Conn) repository.CheckerRepo {
	return &checkerRepoPostgres{conn: conn}
}

const checkExpertExistsQuery = `SELECT user_id FROM experts WHERE expert_id = $1`

func (r *checkerRepoPostgres) CheckExpertExists(expertID string) (bool, error) {
	var userID string

	row := r.conn.QueryRow(context.Background(), checkExpertExistsQuery, expertID)
	err := row.Scan(&userID)

	if err != nil {
		log.Println(err)
		return false, nil
	}
	return true, nil
}
