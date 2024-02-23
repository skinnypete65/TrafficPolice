package models

import "time"

type CaseDTO struct {
	Payload string `json:"payload"`
}

type Case struct {
	ID             string
	Transport      Transport
	Camera         Camera
	Violation      Violation
	ViolationValue string
	RequiredSkill  int
	Date           time.Time
}
