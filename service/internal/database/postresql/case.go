package database

import (
	"TrafficPolice/internal/database"
	"TrafficPolice/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type caseDBPostgres struct {
	conn *pgx.Conn
}

func NewCaseDBPostgres(conn *pgx.Conn) database.CaseDB {
	return &caseDBPostgres{conn: conn}
}

func (db *caseDBPostgres) InsertCase(c *models.Case) error {
	query := `INSERT INTO cases (case_id, transport_id, camera_id, violation_id, violation_value, required_skill, case_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.conn.Exec(context.Background(), query,
		c.ID, c.Transport.ID, c.Camera.ID, c.Violation.ID, c.ViolationValue, c.RequiredSkill, c.Date,
	)
	return err
}
