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

type CameraHandler struct {
	service  services.CameraService
	validate *validator.Validate
}

func NewCameraHandler(service services.CameraService, validate *validator.Validate) *CameraHandler {
	return &CameraHandler{
		service:  service,
		validate: validate,
	}
}

func (h *CameraHandler) AddCameraType(w http.ResponseWriter, r *http.Request) {
	var cameraType domain.CameraType
	err := json.NewDecoder(r.Body).Decode(&cameraType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(cameraType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.AddCameraType(cameraType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("camera type added successfully"))
	if err != nil {
		log.Println(err)
	}
}

func (h *CameraHandler) RegisterCamera(w http.ResponseWriter, r *http.Request) {
	var camera dto.Camera
	err := json.NewDecoder(r.Body).Decode(&camera)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(camera)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.RegisterCamera(domain.Camera{
		ID:         "",
		CameraType: domain.CameraType{ID: camera.CameraTypeID},
		Latitude:   camera.Latitude,
		Longitude:  camera.Longitude,
		ShortDesc:  camera.ShortDesc,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("camera added successfully"))
	if err != nil {
		log.Println(err)
	}
}
