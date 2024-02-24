package models

type Transport struct {
	ID     string
	Chars  string
	Num    string
	Region string
	Person *Person
}
