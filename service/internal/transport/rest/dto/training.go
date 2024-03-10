package dto

type SolvedCasesParams struct {
	CameraID      string `json:"camera_id" validate:"required"`
	RequiredSkill int    `json:"required_skill" validate:"required"`
	ViolationID   string `json:"violation_id" validate:"required"`
	StartTime     string `json:"start_time" validate:"required,is_date_only"`
	EndTime       string `json:"end_time" validate:"required,is_date_only"`
}

type TrainingInfo struct {
	Cases      []Case     `json:"cases"`
	Pagination Pagination `json:"pagination"`
}
