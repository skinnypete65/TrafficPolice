package domain

type RatingInfo struct {
	ExpertID        string
	Username        string
	CompetenceSkill int
	CorrectCnt      int
	IncorrectCnt    int
}

type ExpertRating struct {
	ExpertID     string
	CorrectCnt   int
	IncorrectCnt int
}

type UpdateCompetenceSkill struct {
	ExpertID       string
	ShouldIncrease bool
}
