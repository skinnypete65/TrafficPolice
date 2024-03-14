package dto

type RatingInfo struct {
	ExpertID        string `json:"expert_id"`
	Username        string `json:"username"`
	CompetenceSkill int    `json:"competence_skill"`
	CorrectCnt      int    `json:"correct_cnt"`
	IncorrectCnt    int    `json:"incorrect_cnt"`
}
