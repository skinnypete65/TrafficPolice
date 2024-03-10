package dto

type Transport struct {
	ID     string  `json:"id,omitempty"`
	Chars  string  `json:"chars,omitempty"`
	Num    string  `json:"num,omitempty"`
	Region string  `json:"region,omitempty"`
	Person *Person `json:"person,omitempty"`
}
