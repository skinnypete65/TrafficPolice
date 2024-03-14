package domain

type Decision struct {
	CaseID       string
	Expert       Expert
	FineDecision bool
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
type SolvedCaseDecision struct {
	ExpertID string
	IsRight  bool
}
