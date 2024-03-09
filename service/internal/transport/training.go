package transport

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/dto"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type TrainingHandler struct {
	trainingService services.TrainingService
	validate        *validator.Validate
}

func NewTrainingHandler(trainingService services.TrainingService, validate *validator.Validate) *TrainingHandler {
	return &TrainingHandler{
		trainingService: trainingService,
		validate:        validate,
	}
}

func (h *TrainingHandler) GetSolvedCasesByParams(w http.ResponseWriter, r *http.Request) {
	var params dto.SolvedCasesParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cases, err := h.trainingService.GetSolvedCasesByParams(domain.SolvedCasesParams{
		CameraID:      params.CameraID,
		RequiredSkill: params.RequiredSkill,
		ViolationID:   params.ViolationID,
		StartTime:     params.StartTime,
		EndTime:       params.EndTime,
	})

	dtos := make([]dto.Case, 0, len(cases))
	for _, c := range cases {
		cDto := dto.Case{
			ID: c.ID,
			Transport: dto.Transport{
				ID:     c.Transport.ID,
				Chars:  c.Transport.Chars,
				Num:    c.Transport.Num,
				Region: c.Transport.Region,
			},
			Camera: dto.Camera{
				ID:           c.Camera.ID,
				CameraTypeID: c.Camera.CameraType.ID,
				Latitude:     c.Camera.Latitude,
				Longitude:    c.Camera.Longitude,
				ShortDesc:    c.Camera.ShortDesc,
			},
			Violation: dto.Violation{
				ID:         c.Violation.ID,
				Name:       c.Violation.Name,
				FineAmount: c.Violation.FineAmount,
			},
			ViolationValue: c.ViolationValue,
			RequiredSkill:  c.RequiredSkill,
			IsSolved:       c.IsSolved,
			FineDecision:   c.FineDecision,
			Date:           c.Date,
		}

		dtos = append(dtos, cDto)
	}

	dtoJson, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(dtoJson)
	if err != nil {
		log.Println(err)
	}
}
