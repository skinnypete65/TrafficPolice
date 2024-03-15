package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/transport/rest/dto"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type CameraHandler struct {
	cameraService   service.CameraService
	authService     service.AuthService
	validate        *validator.Validate
	cameraConverter *converter.CameraConverter
}

func NewCameraHandler(
	service service.CameraService,
	authService service.AuthService,
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

// AddCameraType docs
// @Summary Регистрация вида камеры
// @Security ApiKeyAuth
// @Tags camera
// @Description Зарегистрировать новый вид камеры может только директор. Возвращает id вида камеры
// @ID create-camera-type
// @Accept  json
// @Produce  json
// @Param input body dto.CameraTypeIn true "Информация о виде камеры"
// @Success 200 {object} response.IDResponse
// @Failure 400,409 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /camera/type [post]
func (h *CameraHandler) AddCameraType(w http.ResponseWriter, r *http.Request) {
	var cameraType dto.CameraTypeIn
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

// RegisterCamera docs
// @Summary Регистрация камеры
// @Security ApiKeyAuth
// @Tags camera
// @Description Зарегистрировать камеру может только директор. Возвращает id камеры
// @ID create-camera
// @Accept  json
// @Produce  json
// @Param input body dto.RegisterCamera true "Информация о камере"
// @Success 200 {object} response.IDResponse
// @Failure 400,409 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /camera [post]
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
		h.cameraConverter.MapRegisterCameraDtoToDomain(registerInfo.CameraIn, registerInfo.SignUp),
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
