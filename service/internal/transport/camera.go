package transport

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/services"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type CameraHandler struct {
	service services.CameraService
}

func NewCameraHandler(service services.CameraService) *CameraHandler {
	return &CameraHandler{service: service}
}

func (h *CameraHandler) AddCameraType(w http.ResponseWriter, r *http.Request) {
	var cameraType models.CameraType
	err := json.NewDecoder(r.Body).Decode(&cameraType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(cameraType)
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
	var camera models.Camera
	err := json.NewDecoder(r.Body).Decode(&camera)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(camera)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.RegisterCamera(camera)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("camera added successfully"))
	if err != nil {
		log.Println(err)
	}
}
