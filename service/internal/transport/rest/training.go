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
	trainingService      services.TrainingService
	paginationService    services.PaginationService
	validate             *validator.Validate
	caseConverter        *converter.CaseConverter
	paginationConverter  *converter.PaginationConverter
	solvedCasesConverter *converter.SolvedCasesConverter
}

func NewTrainingHandler(
	trainingService services.TrainingService,
	paginationService services.PaginationService,
	validate *validator.Validate,
	caseConverter *converter.CaseConverter,
	paginationConverter *converter.PaginationConverter,
	solvedCasesConverter *converter.SolvedCasesConverter,
) *TrainingHandler {
	return &TrainingHandler{
		trainingService:      trainingService,
		paginationService:    paginationService,
		validate:             validate,
		caseConverter:        caseConverter,
		paginationConverter:  paginationConverter,
		solvedCasesConverter: solvedCasesConverter,
	}
}

func (h *TrainingHandler) GetSolvedCasesByParams(w http.ResponseWriter, r *http.Request) {
	page, err := h.parseQueryParam(r, pageKey, defaultPage)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	limit, err := h.parseQueryParam(r, limitKey, defaultLimit)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	var params dto.SolvedCasesParams

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(params)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	paginationParams := domain.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	cases, err := h.trainingService.GetSolvedCasesByParams(
		h.solvedCasesConverter.MapParamsDtoToDomain(params),
		paginationParams,
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			response.NotFound(w, "Cases with input params are not found")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	pagination, err := h.paginationService.GetPaginationInfo(casesTableName, paginationParams)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	trainingInfo := dto.TrainingInfo{
		Cases:      h.caseConverter.MapDomainsToDto(cases),
		Pagination: h.paginationConverter.MapDomainToDto(pagination),
	}

	infoBody, err := json.Marshal(trainingInfo)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, infoBody)
}

func (h *TrainingHandler) parseQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	queryParam := r.URL.Query().Get(key)

	if queryParam == "" {
		return defaultValue, nil
	}

	page, err := strconv.Atoi(queryParam)
	if err != nil {
		return 0, err
	}

	if page == 0 {
		return defaultValue, nil
	}
	return page, nil

}
