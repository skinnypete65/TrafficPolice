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

func (r *expertRepoPostgres) GetLastNotSolvedCase(expertID string) (string, error) {
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
		return domain.Expert{}, errs.ErrNoLastNotSolvedCase
	}
	if err != nil {
		return domain.Expert{}, err
	}

	return expert, nil
}
