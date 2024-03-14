package dto

type SignInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refresh_token"`
}

type SignUp struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ConfirmExpertInput struct {
	ExpertID    string `json:"expert_id" validate:"required,uuid"`
	IsConfirmed bool   `json:"is_confirmed"`
}
