package dto

type SignInInputDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInOutputDTO struct {
	AccessToken string `json:"accessToken"`
}

type SignUpDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
