package converter

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/dto"
)

type UserInfoConverter struct {
}

func NewUserInfoConverter() *UserInfoConverter {
	return &UserInfoConverter{}
}

func (c *UserInfoConverter) MapSignUpToUserInfo(signUp dto.SignUp) domain.UserInfo {
	return domain.UserInfo{
		Username: signUp.Username,
		Password: signUp.Password,
	}
}

func (c *UserInfoConverter) MapSignInToUserInfo(signIn dto.SignInInput) domain.UserInfo {
	return domain.UserInfo{
		Username: signIn.Username,
		Password: signIn.Password,
	}
}
