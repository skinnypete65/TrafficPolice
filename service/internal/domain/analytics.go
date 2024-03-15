package domain

import "time"

type IntervalCase struct {
	GotAt              time.Time
	IsExpertSolve      bool
	ExpertFineDecision bool
	CaseFineDecision   bool
}

type AnalyticsInterval struct {
	Date                 Date
	AllCases             int
	CorrectCnt           int
	IncorrectCnt         int
	UnknownCnt           int
	MaxConsecutiveSolved int
}
