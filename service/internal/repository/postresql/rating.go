package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
)

type ratingRepoPostgres struct {
	conn *pgx.Conn
}

func NewRatingRepoPostgres(conn *pgx.Conn) repository.RatingRepo {
	return &ratingRepoPostgres{
		conn: conn,
	}
}

const getSolvedCaseDecisionsQuery = `SELECT sc.expert_id, sc.fine_decision = c.fine_decision AS is_right
FROM cases AS c
JOIN solved_cases AS sc ON c.case_id = sc.case_id
WHERE sc.case_id = $1`

func (r *ratingRepoPostgres) GetSolvedCaseDecisions(
	caseDecision domain.CaseDecisionInfo,
) ([]domain.SolvedCaseDecision, error) {
	rows, err := r.conn.Query(context.Background(), getSolvedCaseDecisionsQuery, caseDecision.CaseID)
	if err != nil {
		return nil, err
	}

	solvedDecisions := make([]domain.SolvedCaseDecision, 0)
	for rows.Next() {
		d := domain.SolvedCaseDecision{}
		err = rows.Scan(&d.ExpertID, &d.IsRight)
		if err != nil {
			continue
		}

		solvedDecisions = append(solvedDecisions, d)
	}

	return solvedDecisions, nil
}

const updateCorrectCntQuery = `UPDATE rating
SET correct_cnt = correct_cnt+1
WHERE expert_id = $1`
const updateInCorrectCntQuery = `UPDATE rating
SET incorrect_cnt = incorrect_cnt+1
WHERE expert_id = $1`

func (r *ratingRepoPostgres) SetRating(decisions []domain.SolvedCaseDecision) error {
	batch := &pgx.Batch{}

	for _, d := range decisions {
		if d.IsRight {
			batch.Queue(updateCorrectCntQuery, d.ExpertID)
		} else {
			batch.Queue(updateInCorrectCntQuery, d.ExpertID)
		}
	}

	return r.conn.SendBatch(context.Background(), batch).Close()
}

const insertExpertIdQuery = `INSERT INTO rating (expert_id, correct_cnt, incorrect_cnt)
VALUES ($1, 0, 0) ON CONFLICT DO NOTHING`

func (r *ratingRepoPostgres) InsertExpertId(expertID string) error {
	_, err := r.conn.Exec(context.Background(), insertExpertIdQuery, expertID)
	return err
}

const getRatingQuery = `SELECT r.expert_id, u.username, e.competence_skill,
       r.correct_cnt, r.incorrect_cnt
FROM rating AS r 
JOIN experts as e ON r.expert_id = e.expert_id
JOIN users AS u ON e.user_id = u.user_id`

func (r *ratingRepoPostgres) GetRating() ([]domain.RatingInfo, error) {
	rows, err := r.conn.Query(context.Background(), getRatingQuery)
	if err != nil {
		return nil, err
	}

	infos := make([]domain.RatingInfo, 0)
	for rows.Next() {
		info := domain.RatingInfo{}
		err = rows.Scan(&info.ExpertID, &info.Username, &info.CompetenceSkill, &info.CorrectCnt, &info.IncorrectCnt)

		if err != nil {
			continue
		}

		infos = append(infos, info)
	}

	return infos, nil
}
