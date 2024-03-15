package dto

type AnalyticsInterval struct {
	Date                 string `json:"date"`
	AllCases             int    `json:"all_cases_cnt"`
	CorrectCnt           int    `json:"correct_cnt"`
	IncorrectCnt         int    `json:"incorrect_cnt"`
	UnknownCnt           int    `json:"unknown_cnt"`
	MaxConsecutiveSolved int    `json:"max_consecutive_solved"`
}
