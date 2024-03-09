package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type trainingRepoPostgres struct {
	conn *pgx.Conn
}

func NewTrainingRepoPostgres(conn *pgx.Conn) repository.TrainingRepo {
	return &trainingRepoPostgres{
		conn: conn,
	}
}

const getSolvedCasesByParamsQuery = `SELECT c.case_id, c.violation_value, c.required_skill, 
	c.case_date, c.fine_decision, t.transport_id, t.transport_chars, t.transport_nums,
	t.region, cam.camera_id, type.camera_type_id, type.camera_type_name, 
	cam.camera_latitude, cam.camera_longitude, cam.short_desc,
	v.violation_id, v.violation_name, v.fine_amount
FROM cases AS c
JOIN transports AS t ON c.transport_id = t.transport_id 
JOIN cameras AS cam ON c.camera_id = cam.camera_id
JOIN camera_types AS type ON cam.camera_type_id = type.camera_type_id
JOIN violations AS v ON c.violation_id = v.violation_id
WHERE c.is_solved = true AND cam.camera_id = $1 AND c.required_skill = $2
AND v.violation_id = $3
AND c.case_date BETWEEN $4 and $5`

func (r *trainingRepoPostgres) GetSolvedCasesByParams(params domain.SolvedCasesParams) ([]domain.Case, error) {
	rows, err := r.conn.Query(context.Background(), getSolvedCasesByParamsQuery,
		params.CameraID, params.RequiredSkill, params.ViolationID, params.StartTime, params.EndTime)
	if err != nil {
		return nil, err
	}

	cases := make([]domain.Case, 0)
	for rows.Next() {
		var c domain.Case

		err = rows.Scan(&c.ID, &c.ViolationValue, &c.RequiredSkill, &c.Date, &c.FineDecision,
			&c.Transport.ID, &c.Transport.Chars, &c.Transport.Num,
			&c.Transport.Region, &c.Camera.ID, &c.Camera.CameraType.ID, &c.Camera.CameraType.Name,
			&c.Camera.Latitude, &c.Camera.Longitude, &c.Violation.ID,
			&c.Violation.ID, &c.Violation.Name, &c.Violation.FineAmount)

		if err != nil {
			log.Println(err)
			continue
		}

		cases = append(cases, c)
	}

	return cases, nil
}
