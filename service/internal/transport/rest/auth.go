package rest

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/dto"
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.RegisterExpert(domain.UserInfo{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input dto.SignInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.SignIn(domain.UserInfo{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ansBytes, err := json.Marshal(dto.SignInOutput{AccessToken: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(ansBytes)
	if err != nil {
		log.Println(err)
	}
}

func (h *AuthHandler) ConfirmExpert(w http.ResponseWriter, r *http.Request) {
	var input dto.ConfirmExpertInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.ConfirmExpert(domain.ConfirmExpert{
		ExpertID:    input.ExpertID,
		IsConfirmed: input.IsConfirmed,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
