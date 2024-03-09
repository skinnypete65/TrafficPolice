package repository

import (
	"TrafficPolice/internal/repository"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type paginationRepoPostgres struct {
	conn *pgx.Conn
}

func NewPaginationRepoPostgres(conn *pgx.Conn) repository.PaginationRepo {
	return &paginationRepoPostgres{
		conn: conn,
	}
}

func (r *paginationRepoPostgres) GetRecordsCount(table string) (int, error) {
	sqlTableQuery := fmt.Sprintf("SELECT count(*) FROM %s", table)
	row := r.conn.QueryRow(context.Background(), sqlTableQuery)

	var recordsCount int
	err := row.Scan(&recordsCount)
	if err != nil {
		return 0, err
	}

	return recordsCount, nil
}
