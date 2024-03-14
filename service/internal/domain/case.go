package domain

import (
	"time"
)

type CaseDTO struct {
	Payload string `json:"payload"`
}

type Case struct {
	ID             string
	Transport      Transport
	Camera         Camera
	Violation      Violation
	ViolationValue string
	RequiredSkill  int64
	Date           time.Time
	IsSolved       bool
	FineDecision   bool
}

type CaseAssessment struct {
	ExpertID      string
	IsExpertSolve bool
	FineDecision  bool
}

type CaseStatus struct {
	CaseID          string
	ViolationValue  string
	RequiredSkill   int
	CaseDate        time.Time
	IsSolved        bool
	FineDecision    bool
	CaseAssessments []CaseAssessment
}
