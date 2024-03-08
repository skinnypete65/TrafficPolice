package domain

type SolvedCase struct {
	SolvedCaseID  string
	ExpertID      string
	CaseID        string
	IsExpertSolve bool
	FineDecision  bool
}
