package dto

type Violation struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	FineAmount int    `json:"fine_amount,omitempty"`
}
