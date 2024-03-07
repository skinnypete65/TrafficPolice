package dto

type Decision struct {
	CaseID       string `json:"case_id,omitempty"`
	FineDecision bool   `json:"fine_decision,omitempty"`
}
