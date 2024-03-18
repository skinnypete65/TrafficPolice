package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type directorRepoPostgres struct {
	conn *pgx.Conn
}

func NewDirectorRepoPostgres(conn *pgx.Conn) repository.DirectorRepo {
	return &directorRepoPostgres{conn: conn}
}

const getCaseQuery = `SELECT case_id, violation_value, required_skill, case_date, is_solved, 
       fine_decision, solved_at
FROM cases WHERE case_id = $1`

const getCaseAssessments = `SELECT expert_id, is_expert_solve, fine_decision 
FROM expert_cases WHERE case_id = $1`

func (r *directorRepoPostgres) GetCase(caseID string) (domain.CaseStatus, error) {
	row := r.conn.QueryRow(context.Background(), getCaseQuery, caseID)

	status := domain.CaseStatus{CaseAssessments: make([]domain.CaseAssessment, 0)}
	err := row.Scan(&status.CaseID, &status.ViolationValue, &status.RequiredSkill, &status.CaseDate,
		&status.IsSolved, &status.FineDecision, &status.SolvedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CaseStatus{}, errs.ErrNoCase
		}
		return domain.CaseStatus{}, err
	}

	assessmentsRows, err := r.conn.Query(context.Background(), getCaseAssessments, caseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, nil
		}
		return domain.CaseStatus{}, err
	}

	assessments := make([]domain.CaseAssessment, 0)
	for assessmentsRows.Next() {
		assessment := domain.CaseAssessment{}
		err = assessmentsRows.Scan(
			&assessment.ExpertID, &assessment.IsExpertSolve, &assessment.FineDecision,
		)
		if err != nil {
			log.Println(err)
			continue
		}

		assessments = append(assessments, assessment)
	}

	status.CaseAssessments = assessments

	return status, nil
}

const getExpertIntervalCasesQuery = `SELECT ec.is_expert_solve , ec.fine_decision AS expert_fine_decision,
c.fine_decision AS case_fine_decision, ec.got_at
FROM expert_cases AS ec
JOIN cases AS c ON ec.case_id = c.case_id
WHERE expert_id = $1
AND (ec.got_at BETWEEN $2 AND $3)
ORDER BY ec.got_at`

func (r *directorRepoPostgres) GetExpertIntervalCases(
	expertID string,
	startDate time.Time,
	endDate time.Time,
) (map[domain.Date][]domain.IntervalCase, error) {
	rows, err := r.conn.Query(context.Background(), getExpertIntervalCasesQuery, expertID, startDate, endDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNoRows
		}
	}

	intervals := make(map[domain.Date][]domain.IntervalCase)
	for rows.Next() {
		interval := domain.IntervalCase{}

		err = rows.Scan(&interval.IsExpertSolve, &interval.ExpertFineDecision, &interval.CaseFineDecision, &interval.GotAt)
		if err != nil {
			log.Println(err)
			continue
		}

		date := domain.NewDate(interval.GotAt.Date())

		if _, ok := intervals[date]; !ok {
			intervals[date] = make([]domain.IntervalCase, 0)
		}
		intervals[date] = append(intervals[date], interval)
	}

	return intervals, nil
}
