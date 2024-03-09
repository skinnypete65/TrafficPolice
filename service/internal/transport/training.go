package transport

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/dto"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultPage    = 1
	defaultLimit   = 10
	pageKey        = "page"
	limitKey       = "limit"
	casesTableName = "cases"
)

type TrainingHandler struct {
	trainingService   services.TrainingService
	paginationService services.PaginationService
	validate          *validator.Validate
}

func NewTrainingHandler(
	trainingService services.TrainingService,
	paginationService services.PaginationService,
	validate *validator.Validate,
) *TrainingHandler {
	return &TrainingHandler{
		trainingService:   trainingService,
		paginationService: paginationService,
		validate:          validate,
	}
}

func (h *TrainingHandler) GetSolvedCasesByParams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, err := strconv.Atoi(r.URL.Query().Get(pageKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if page == 0 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(r.URL.Query().Get(limitKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if limit == 0 {
		limit = defaultLimit
	}

	var params dto.SolvedCasesParams

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paginationParams := domain.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	cases, err := h.trainingService.GetSolvedCasesByParams(
		domain.SolvedCasesParams{
			CameraID:      params.CameraID,
			RequiredSkill: params.RequiredSkill,
			ViolationID:   params.ViolationID,
			StartTime:     params.StartTime,
			EndTime:       params.EndTime,
		},
		paginationParams,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pagination, err := h.paginationService.GetPaginationInfo(casesTableName, paginationParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trainingInfo := dto.TrainingInfo{
		Cases:      mapCasesToDTO(cases),
		Pagination: mapPaginationToDTO(pagination),
	}

	infoJson, err := json.Marshal(trainingInfo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(infoJson)
	if err != nil {
		log.Println(err)
	}
}

func mapCasesToDTO(cases []domain.Case) []dto.Case {
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

	return dtos
}

func mapPaginationToDTO(pagination domain.Pagination) dto.Pagination {
	return dto.Pagination{
		Next:          pagination.Next,
		Previous:      pagination.Next,
		RecordPerPage: pagination.RecordPerPage,
		CurrentPage:   pagination.CurrentPage,
		TotalPage:     pagination.TotalPage,
	}
}
