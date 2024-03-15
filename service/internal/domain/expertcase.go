package domain

import "time"

type ExpertCase struct {
	ExpertCaseID  string
	ExpertID      string
	CaseID        string
	IsExpertSolve bool
	FineDecision  bool
	GotAt         time.Time
	SolvedAt      time.Time
}
