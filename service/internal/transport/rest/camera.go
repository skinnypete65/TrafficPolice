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

type CameraHandler struct {
	cameraService services.CameraService
	authService   services.AuthService
	validate      *validator.Validate
}

func NewCameraHandler(
	service services.CameraService,
	authService services.AuthService,
	validate *validator.Validate,
) *CameraHandler {
	return &CameraHandler{
		cameraService: service,
		authService:   authService,
		validate:      validate,
	}
}

func (h *CameraHandler) AddCameraType(w http.ResponseWriter, r *http.Request) {
	var cameraType dto.CameraType
	err := json.NewDecoder(r.Body).Decode(&cameraType)

	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(cameraType)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	cameraTypeID, err := h.cameraService.AddCameraType(domain.CameraType{
		Name: cameraType.Name,
	})

	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			conflict(w, "Camera with input name already exists")
			return
		}
		internalServerError(w)
		return
	}

	idResponse(w, cameraTypeID)
}

func (h *CameraHandler) RegisterCamera(w http.ResponseWriter, r *http.Request) {
	var registerInfo dto.RegisterCamera
	err := json.NewDecoder(r.Body).Decode(&registerInfo)

	if err != nil {
		badRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(registerInfo)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	cameraID, err := h.authService.RegisterCamera(domain.RegisterCamera{
		Camera: domain.Camera{
			ID:         "",
			CameraType: domain.CameraType{ID: registerInfo.Camera.CameraTypeID},
			Latitude:   registerInfo.Camera.Latitude,
			Longitude:  registerInfo.Camera.Longitude,
			ShortDesc:  registerInfo.Camera.ShortDesc,
		},
		Username: registerInfo.SignUp.Username,
		Password: registerInfo.SignUp.Password,
	})
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			conflict(w, "Camera with this username already exists")
			return
		}
		log.Println(err)
		internalServerError(w)
		return
	}

	idResponse(w, cameraID)
}
