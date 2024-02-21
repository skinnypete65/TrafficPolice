package transport

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/services/service"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type CameraHandler struct {
	service service.CameraService
}

func NewCameraHandler(service service.CameraService) *CameraHandler {
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

	w.Write([]byte("camera type added successfully"))
}
