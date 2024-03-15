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
	trainingService      service.TrainingService
	paginationService    service.PaginationService
	validate             *validator.Validate
	caseConverter        *converter.CaseConverter
	paginationConverter  *converter.PaginationConverter
	solvedCasesConverter *converter.SolvedCasesConverter
}

func NewTrainingHandler(
	trainingService service.TrainingService,
	paginationService service.PaginationService,
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

// GetSolvedCasesByParams docs
// @Summary Получение проишествий для тренировки
// @Security ApiKeyAuth
// @Tags expert
// @Description Получение прошествий для тренировки. Может воспользоваться только эксперт
// @ID expert-training
// @Accept  json
// @Produce  json
// @Param input body dto.SolvedCasesParams true "Информация для фильтров по проишествиям"
// @Param page query int true "номер страницы"
// @Param limit query int true "Лимит кейсов на странице"
// @Success 200 {object} dto.TrainingInfo
// @Failure 400,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /expert/training [post]
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

	param, err := strconv.Atoi(queryParam)
	if err != nil {
		return 0, err
	}

	if param == 0 {
		return defaultValue, nil
	}
	return param, nil

}
