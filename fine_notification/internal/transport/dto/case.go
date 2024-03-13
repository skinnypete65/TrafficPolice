package dto

import "time"

type Case struct {
	ID             string    `json:"id,omitempty"`
	Transport      Transport `json:"middlewares,omitempty"`
	Camera         Camera    `json:"camera,omitempty"`
	Violation      Violation `json:"violation,omitempty"`
	ViolationValue string    `json:"violation_value,omitempty"`
	RequiredSkill  int       `json:"required_skill,omitempty"`
	Date           time.Time `json:"date,omitempty"`
	IsSolved       bool      `json:"is_solved,omitempty"`
	FineDecision   bool      `json:"fine_decision,omitempty"`
}

type CaseWithImage struct {
	Case           Case   `json:"case"`
	Image          []byte `json:"image"`
	ImageExtension string `json:"image_extension"`
}
