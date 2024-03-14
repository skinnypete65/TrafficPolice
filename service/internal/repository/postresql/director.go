package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type directorRepoPostgres struct {
	conn *pgx.Conn
}

func NewDirectorRepoPostgres(conn *pgx.Conn) repository.DirectorRepo {
	return &directorRepoPostgres{conn: conn}
}

const getCasesQuery = `SELECT case_id, violation_value, required_skill, case_date, is_solved, fine_decision
FROM cases`

const getCaseAssessments = `SELECT expert_id, is_expert_solve, fine_decision 
FROM solved_cases WHERE case_id = $1`

func (r *directorRepoPostgres) GetCases() ([]domain.CaseStatus, error) {
	rows, err := r.conn.Query(context.Background(), getCasesQuery)
	if err != nil {
		return nil, err
	}

	statuses := make([]domain.CaseStatus, 0)
	for rows.Next() {
		status := domain.CaseStatus{CaseAssessments: make([]domain.CaseAssessment, 0)}

		err = rows.Scan(&status.CaseID, &status.ViolationValue, &status.RequiredSkill, &status.CaseDate, &status.IsSolved, &status.FineDecision)
		if err != nil {
			log.Println(err)
			continue
		}

		statuses = append(statuses, status)
	}

	for i := range statuses {
		assessmentsRows, err := r.conn.Query(context.Background(), getCaseAssessments, statuses[i].CaseID)
		if err != nil {
			log.Println(err)
			continue
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

		statuses[i].CaseAssessments = assessments
	}

	return statuses, nil
}
