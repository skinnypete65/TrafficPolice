package repository

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type caseRepoPostgres struct {
	conn *pgx.Conn
}

func NewCaseRepoPostgres(conn *pgx.Conn) repository.CaseRepo {
	return &caseRepoPostgres{conn: conn}
}

func (r *caseRepoPostgres) InsertCase(c *domain.Case) error {
	query := `INSERT INTO cases (case_id, transport_id, camera_id, 
                   violation_id, violation_value, required_skill, case_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.conn.Exec(context.Background(), query,
		c.ID, c.Transport.ID, c.Camera.ID, c.Violation.ID, c.ViolationValue, c.RequiredSkill, c.Date,
	)
	return err
}

const getCaseByIDQuery = `SELECT c.case_id, t.transport_id, t.transport_chars, 
       t.transport_nums, t.region, t.person_id, cam.camera_type_id,
       cam.camera_latitude, cam.camera_longitude, cam.short_desc, v.violation_name, v.fine_amount,
       c.violation_value, c.required_skill, c.case_date,
       c.is_solved, c.fine_decision
FROM cases as c
JOIN transports AS t ON c.transport_id = t.transport_id
JOIN violations AS v ON c.violation_id = v.violation_id
JOIN cameras AS cam ON c.camera_id = cam.camera_id
WHERE c.case_id = $1
LIMIT 1`

func (r *caseRepoPostgres) GetCaseByID(caseID string) (domain.Case, error) {
	c := domain.Case{Transport: domain.Transport{Person: &domain.Person{}}, Camera: domain.Camera{}, Violation: domain.Violation{}}

	row := r.conn.QueryRow(context.Background(), getCaseByIDQuery, caseID)

	err := row.Scan(&c.ID, &c.Transport.ID, &c.Transport.Chars, &c.Transport.Num, &c.Transport.Region,
		&c.Transport.Person.ID, &c.Camera.CameraTypeID, &c.Camera.Latitude, &c.Camera.Longitude,
		&c.Camera.ShortDesc, &c.Violation.Name, &c.Violation.FineAmount, &c.ViolationValue,
		&c.RequiredSkill, &c.Date, &c.IsSolved, &c.FineDecision)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Case{}, errs.ErrNoCase
	}
	if err != nil {
		return domain.Case{}, err
	}

	return c, nil
}
