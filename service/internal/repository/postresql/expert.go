package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
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

const getLastNotSolvedCaseQuery = `SELECT case_id FROM expert_cases 
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

const getNotSolvedCaseQuery = `SELECT c.case_id, t.transport_id, t.transport_chars, 
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
	SELECT ec.case_id FROM expert_cases as ec
	WHERE ec.expert_id = $2 and ec.is_expert_solve = true
)
LIMIT 1`

func (r *expertRepoPostgres) GetNotSolvedCase(expert domain.Expert) (domain.Case, error) {
	c := domain.Case{Transport: domain.Transport{Person: &domain.Person{}},
		Camera: domain.Camera{}, Violation: domain.Violation{},
	}

	row := r.conn.QueryRow(context.Background(), getNotSolvedCaseQuery, expert.CompetenceSkill, expert.ID)

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

const insertNotSolvedCaseQuery = `INSERT INTO expert_cases 
    (expert_case_id, expert_id, case_id, is_expert_solve, fine_decision, got_at, solved_at) 
    VALUES ($1, $2, $3, $4, $5, $6, NULL)`

func (r *expertRepoPostgres) InsertNotSolvedCase(solvedCase domain.ExpertCase) error {
	_, err := r.conn.Exec(context.Background(), insertNotSolvedCaseQuery,
		solvedCase.ExpertCaseID,
		solvedCase.ExpertID,
		solvedCase.CaseID,
		solvedCase.IsExpertSolve,
		solvedCase.FineDecision,
		solvedCase.GotAt,
	)

	return err
}

const setCaseDecisionQuery = `UPDATE expert_cases 
SET is_expert_solve = true, fine_decision = $1, solved_at = $2
WHERE expert_id = $3 and case_id = $4`

func (r *expertRepoPostgres) SetCaseDecision(decision domain.Decision) error {
	_, err := r.conn.Exec(context.Background(), setCaseDecisionQuery,
		decision.FineDecision,
		decision.SolvedAt,
		decision.Expert.ID,
		decision.CaseID,
	)

	return err
}

const gGetCaseFineDecisions = `SELECT 
    SUM(CASE WHEN ec.fine_decision = true THEN 1 ELSE 0 END) AS positive_decisions,
    SUM(CASE WHEN ec.fine_decision = false THEN 1 ELSE 0 END) AS negative_decisions
FROM expert_cases AS ec
JOIN experts AS e ON ec.expert_id = e.expert_id 
WHERE ec.case_id = $1 and ec.is_expert_solve = true and e.competence_skill = $2`

func (r *expertRepoPostgres) GetCaseFineDecisions(caseID string, competenceSkill int) (domain.FineDecisions, error) {
	row := r.conn.QueryRow(context.Background(), gGetCaseFineDecisions, caseID, competenceSkill)

	var fineDecisions domain.FineDecisions
	err := row.Scan(&fineDecisions.PositiveDecisions, &fineDecisions.NegativeDecisions)
	return fineDecisions, err
}

const getExpertsCountByRequiredSkillQuery = `SELECT COUNT(*)
FROM experts
WHERE competence_skill = $1 and is_confirmed = true`

func (r *expertRepoPostgres) GetExpertsCountBySkill(competenceSkill int) (int, error) {
	row := r.conn.QueryRow(context.Background(), getExpertsCountByRequiredSkillQuery, competenceSkill)

	var cnt int
	err := row.Scan(&cnt)
	return cnt, err
}
