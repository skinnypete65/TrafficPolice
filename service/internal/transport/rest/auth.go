package rest

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/dto"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type AuthHandler struct {
	service  services.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service services.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validate,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input dto.SignUp

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.service.RegisterExpert(domain.UserInfo{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			conflict(w, "User with this username already exists")
			return
		}
		log.Println(err)
		internalServerError(w)
		return
	}
	oKMessage(w, "You signed up successfully")
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input dto.SignInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	token, err := h.service.SignIn(domain.UserInfo{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			notFound(w, "User with this username not exists")
			return
		}
		if errors.Is(err, errs.ErrInvalidPass) {
			unauthorized(w)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ansBytes, err := json.Marshal(dto.SignInOutput{AccessToken: token})
	if err != nil {
		internalServerError(w)
		log.Println(err)
		return
	}

	writeResponse(w, http.StatusOK, ansBytes)
	if err != nil {
		log.Println(err)
	}
}

func (h *AuthHandler) ConfirmExpert(w http.ResponseWriter, r *http.Request) {
	var input dto.ConfirmExpertInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.service.ConfirmExpert(domain.ConfirmExpert{
		ExpertID:    input.ExpertID,
		IsConfirmed: input.IsConfirmed,
	})

	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			notFound(w, "Expert with input id not found")
			return
		}
		internalServerError(w)
		return
	}

	if input.IsConfirmed {
		oKMessage(w, "Expert confirmed successfully")
	} else {
		oKMessage(w, "Expert unconfirmed successfully")
	}
}
