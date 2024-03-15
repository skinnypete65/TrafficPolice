package domain

import "time"

type Decision struct {
	CaseID       string
	Expert       Expert
	FineDecision bool
	SolvedAt     time.Time
}

type FineDecisions struct {
	PositiveDecisions int
	NegativeDecisions int
}

type CaseDecisionInfo struct {
	CaseID         string
	ShouldSendFine bool
	IsSolved       bool
}
type ExpertCaseDecision struct {
	ExpertID string
	IsRight  bool
}
