package rest

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/dto"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type CameraHandler struct {
	cameraService   services.CameraService
	authService     services.AuthService
	validate        *validator.Validate
	cameraConverter *converter.CameraConverter
}

func NewCameraHandler(
	service services.CameraService,
	authService services.AuthService,
	validate *validator.Validate,
	cameraConverter *converter.CameraConverter,
) *CameraHandler {
	return &CameraHandler{
		cameraService:   service,
		authService:     authService,
		validate:        validate,
		cameraConverter: cameraConverter,
	}
}

func (h *CameraHandler) AddCameraType(w http.ResponseWriter, r *http.Request) {
	var cameraType dto.CameraType
	err := json.NewDecoder(r.Body).Decode(&cameraType)

	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(cameraType)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	cameraTypeID, err := h.cameraService.AddCameraType(domain.CameraType{
		Name: cameraType.Name,
	})

	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			response.Conflict(w, "Camera with input name already exists")
			return
		}
		response.InternalServerError(w)
		return
	}

	response.IdResponse(w, cameraTypeID)
}

func (h *CameraHandler) RegisterCamera(w http.ResponseWriter, r *http.Request) {
	var registerInfo dto.RegisterCamera
	err := json.NewDecoder(r.Body).Decode(&registerInfo)

	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(registerInfo)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	cameraID, err := h.authService.RegisterCamera(
		h.cameraConverter.MapRegisterCameraDomainToDto(registerInfo.Camera, registerInfo.SignUp),
	)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			response.Conflict(w, "Camera with this username already exists")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.IdResponse(w, cameraID)
}
