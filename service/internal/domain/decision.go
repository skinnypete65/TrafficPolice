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
