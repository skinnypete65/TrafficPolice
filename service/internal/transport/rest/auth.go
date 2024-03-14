package rest

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/dto"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type AuthHandler struct {
	service           services.AuthService
	validate          *validator.Validate
	userInfoConverter *converter.UserInfoConverter
	authConverter     *converter.AuthConverter
}

func NewAuthHandler(service services.AuthService,
	validate *validator.Validate,
	userInfoConverter *converter.UserInfoConverter,
	authConverter *converter.AuthConverter,
) *AuthHandler {
	return &AuthHandler{
		service:           service,
		validate:          validate,
		userInfoConverter: userInfoConverter,
		authConverter:     authConverter,
	}
}

// SignUp docs
// @Summary Регистрация эксперта
// @Tags auth
// @Description Регистрация эксперта по логину и паролю
// @ID auth-sign-up
// @Accept  json
// @Produce  json
// @Param input body dto.SignUp true "Логин и пароль"
// @Success 200 {object} response.Body
// @Failure 400,401,409 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /auth/sign_up [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input dto.SignUp

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.service.RegisterExpert(h.userInfoConverter.MapSignUpToUserInfo(input))
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			response.Conflict(w, "User with this username already exists")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}
	response.OKMessage(w, "You signed up successfully")
}

// SignIn docs
// @Summary Вход пользователей
// @Tags auth
// @Description Вход для всех пользователей по логину и паролю
// @ID auth-sign-in
// @Accept  json
// @Produce  json
// @Param input body dto.SignInInput true "Логин и пароль"
// @Success 200 {object} dto.SignInOutput
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /auth/sign_in [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input dto.SignInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	tokens, err := h.service.SignIn(h.userInfoConverter.MapSignInToUserInfo(input))

	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			response.NotFound(w, "User with this username not exists")
			return
		}
		if errors.Is(err, errs.ErrInvalidPass) {
			response.Unauthorized(w)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ansBytes, err := json.Marshal(
		dto.SignInOutput{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	)
	if err != nil {
		response.InternalServerError(w)
		log.Println(err)
		return
	}

	response.WriteResponse(w, http.StatusOK, ansBytes)
	if err != nil {
		log.Println(err)
	}
}

// ConfirmExpert docs
// @Summary Подтверждение эксперта
// @Security ApiKeyAuth
// @Tags auth
// @Description Подтверждение эксперта по id. Может делать только директор
// @ID auth-confirm-expert
// @Accept  json
// @Produce  json
// @Param input body dto.ConfirmExpertInput true "id эксперта и информация о подтверждении"
// @Success 200 {object} dto.SignInOutput
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /auth/confirm/expert [post]
func (h *AuthHandler) ConfirmExpert(w http.ResponseWriter, r *http.Request) {
	var input dto.ConfirmExpertInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.service.ConfirmExpert(h.authConverter.MapConfirmExpertDtoToDomain(input))

	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			response.NotFound(w, "Expert with input id not found")
			return
		}
		response.InternalServerError(w)
		return
	}

	if input.IsConfirmed {
		response.OKMessage(w, "Expert confirmed successfully")
	} else {
		response.OKMessage(w, "Expert unconfirmed successfully")
	}
}
