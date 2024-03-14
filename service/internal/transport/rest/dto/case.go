package dto

import "time"

type Case struct {
	ID             string    `json:"id,omitempty"`
	Transport      Transport `json:"transport,omitempty"`
	Camera         Camera    `json:"camera,omitempty"`
	Violation      Violation `json:"violation,omitempty"`
	ViolationValue string    `json:"violation_value,omitempty"`
	RequiredSkill  int64     `json:"required_skill,omitempty"`
	Date           time.Time `json:"date,omitempty"`
	IsSolved       bool      `json:"is_solved,omitempty"`
	FineDecision   bool      `json:"fine_decision,omitempty"`
}

type CaseWithImage struct {
	Case           Case   `json:"case"`
	Image          []byte `json:"image"`
	ImageExtension string `json:"image_extension"`
}

type CaseAssessment struct {
	ExpertID      string `json:"expert_id"`
	IsExpertSolve bool   `json:"is_expert_solve"`
	FineDecision  bool   `json:"fine_decision"`
}

type CaseStatus struct {
	CaseID          string           `json:"case_id"`
	ViolationValue  string           `json:"violation_value"`
	RequiredSkill   int              `json:"required_skill"`
	CaseDate        time.Time        `json:"case_date"`
	IsSolved        bool             `json:"is_solved"`
	FineDecision    bool             `json:"fine_decision"`
	CaseAssessments []CaseAssessment `json:"case_assessments"`
}
