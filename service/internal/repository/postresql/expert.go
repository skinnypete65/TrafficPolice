package repository

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type expertRepoPostgres struct {
	conn *pgx.Conn
}

func NewExpertRepoPostgres(conn *pgx.Conn) repository.ExpertRepo {
	return &expertRepoPostgres{conn: conn}
}

const getLastNotSolvedCaseQuery = `SELECT case_id FROM solved_cases 
	WHERE expert_id = $1 and is_expert_solve = false`

func (r *expertRepoPostgres) GetLastNotSolvedCaseID(expertID string) (string, error) {
	var caseID string

	row := r.conn.QueryRow(context.Background(), getLastNotSolvedCaseQuery, expertID)
	err := row.Scan(&caseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errs.ErrNoLastNotSolvedCase
	}
	if err != nil {
		return "", err
	}

	return caseID, err
}

const getExpertByUserIDQuery = `SELECT e.expert_id, e.is_confirmed, e.competence_skill,
       u.user_id, u.username, u.hash_pass, u.register_at, u.role
	FROM experts AS e
	JOIN users AS u on e.user_id = u.user_id
	WHERE u.user_id = $1`

func (r *expertRepoPostgres) GetExpertByUserID(userID string) (domain.Expert, error) {
	expert := domain.Expert{UserInfo: domain.UserInfo{}}

	row := r.conn.QueryRow(context.Background(), getExpertByUserIDQuery, userID)

	err := row.Scan(&expert.ID, &expert.IsConfirmed, &expert.CompetenceSkill,
		&expert.UserInfo.ID, &expert.UserInfo.Username, &expert.UserInfo.Password,
		&expert.UserInfo.RegisterAt, &expert.UserInfo.UserRole)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Expert{}, errs.ErrUserNotExists
	}
	if err != nil {
		return domain.Expert{}, err
	}

	return expert, nil
}

const getNotSolvedCasesQuery = `SELECT c.case_id, t.transport_id, t.transport_chars, 
       t.transport_nums, t.region, t.person_id, cam.camera_type_id,
       cam.camera_latitude, cam.camera_longitude, cam.short_desc, v.violation_name, v.fine_amount,
       c.violation_value, c.required_skill, c.case_date,
       c.is_solved, c.fine_decision
FROM cases as c
JOIN transports AS t ON c.transport_id = t.transport_id
JOIN violations AS v ON c.violation_id = v.violation_id
JOIN cameras AS cam ON c.camera_id = cam.camera_id
WHERE c.is_solved = false and c.required_skill = $1 and c.case_id 
NOT IN (
	SELECT sc.case_id FROM solved_cases as sc
	WHERE sc.expert_id = $2 and sc.is_expert_solve = true
)
LIMIT 1`

func (r *expertRepoPostgres) GetNotSolvedCase(expert domain.Expert) (domain.Case, error) {
	c := domain.Case{Transport: domain.Transport{Person: &domain.Person{}},
		Camera: domain.Camera{}, Violation: domain.Violation{},
	}

	row := r.conn.QueryRow(context.Background(), getNotSolvedCasesQuery, expert.CompetenceSkill, expert.ID)

	err := row.Scan(&c.ID, &c.Transport.ID, &c.Transport.Chars, &c.Transport.Num, &c.Transport.Region,
		&c.Transport.Person.ID, &c.Camera.CameraType.ID, &c.Camera.Latitude, &c.Camera.Longitude,
		&c.Camera.ShortDesc, &c.Violation.Name, &c.Violation.FineAmount, &c.ViolationValue,
		&c.RequiredSkill, &c.Date, &c.IsSolved, &c.FineDecision)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Case{}, errs.ErrNoNotSolvedCase
	}
	if err != nil {
		return domain.Case{}, err
	}

	return c, nil
}

const insertNotSolvedCaseQuery = `INSERT INTO solved_cases 
    (solved_case_id, expert_id, case_id, is_expert_solve, fine_decision) 
    VALUES ($1, $2, $3, $4, $5)`

func (r *expertRepoPostgres) InsertNotSolvedCase(solvedCase domain.SolvedCase) error {
	_, err := r.conn.Exec(context.Background(), insertNotSolvedCaseQuery,
		solvedCase.SolvedCaseID,
		solvedCase.ExpertID,
		solvedCase.CaseID,
		solvedCase.IsExpertSolve,
		solvedCase.FineDecision,
	)

	return err
}

const setCaseDecisionQuery = `UPDATE solved_cases SET is_expert_solve = true, fine_decision = $1
    WHERE expert_id = $2 and case_id = $3`

func (r *expertRepoPostgres) SetCaseDecision(decision domain.Decision) error {
	_, err := r.conn.Exec(context.Background(), setCaseDecisionQuery,
		decision.FineDecision,
		decision.Expert.ID,
		decision.CaseID,
	)

	return err
}

const gGetCaseFineDecisions = `SELECT 
    SUM(CASE WHEN fine_decision = true THEN 1 ELSE 0 END) AS positive_decisions,
    SUM(CASE WHEN fine_decision = false THEN 1 ELSE 0 END) AS negative_decisions
FROM solved_cases
WHERE case_id = $1`

func (r *expertRepoPostgres) GetCaseFineDecisions(caseID string) (domain.FineDecisions, error) {
	row := r.conn.QueryRow(context.Background(), gGetCaseFineDecisions, caseID)

	var fineDecisions domain.FineDecisions
	err := row.Scan(&fineDecisions.PositiveDecisions, &fineDecisions.NegativeDecisions)
	return fineDecisions, err
}

const getExpertsCountByRequiredSkillQuery = `SELECT COUNT(*)
FROM experts
WHERE competence_skill = $1`

func (r *expertRepoPostgres) GetExpertsCountBySkill(competenceSkill int) (int, error) {
	row := r.conn.QueryRow(context.Background(), getExpertsCountByRequiredSkillQuery, competenceSkill)

	var cnt int
	err := row.Scan(&cnt)
	return cnt, err
}
