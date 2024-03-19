package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/transport/rest/dto"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

const (
	expertIDKey  = "id"
	startTimeKey = "start_time"
	endTimeKey   = "end_time"
	caseIDKey    = "id"
)

type DirectorHandler struct {
	directorService    service.DirectorService
	caseConverter      *converter.CaseConverter
	analyticsConverter *converter.AnalyticsConverter
}

func NewDirectorHandler(
	directorService service.DirectorService,
	caseConverter *converter.CaseConverter,
	analyticsConverter *converter.AnalyticsConverter,
) *DirectorHandler {
	return &DirectorHandler{
		directorService:    directorService,
		caseConverter:      caseConverter,
		analyticsConverter: analyticsConverter,
	}
}

// GetCase docs
// @Summary Получение состояния для случая
// @Security ApiKeyAuth
// @Tags director
// @Description Получение состояния для конкретного случая по его id. Воспользоваться может только директор
// @ID director-case-get
// @Produce  json
// @Param id query string true "id случая"
// @Success 200 {object} dto.CaseStatus
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /director/case [get]
func (h *DirectorHandler) GetCase(w http.ResponseWriter, r *http.Request) {
	caseID := r.URL.Query().Get(caseIDKey)
	if caseID == "" {
		response.BadRequest(w, "id is empty")
		return
	}
	_, err := uuid.Parse(caseID)
	if err != nil {
		response.BadRequest(w, "case id is not uuid")
		return
	}

	caseStatus, err := h.directorService.GetCase(caseID)
	if err != nil {
		if errors.Is(err, errs.ErrNoCase) {
			response.NotFound(w, "Case with input id not found")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	casesBytes, err := json.Marshal(h.caseConverter.MapCaseStatusToDto(caseStatus))
	if err != nil {
		response.InternalServerError(w)
	}

	response.WriteResponse(w, http.StatusOK, casesBytes)
}

// ExpertAnalytics docs
// @Summary Получение аналитики проверяющих специалистов по промежуткам времени
// @Security ApiKeyAuth
// @Tags director
// @Description Получить количество всех случаев, правильно решенных случаев, неправильно решенных случаев, неизвестных случаев и максимальное количество подряд решенных задач. Воспользоваться может только директор
// @ID director-analytics-expert
// @Produce  json
// @Param id query string true "id эксперта"
// @Param start_time query string true "Начало промежутка времени в формате yyyy-mm-dd"
// @Param end_time query string true "Конец промежутка времени в формате yyyy-mm-dd"
// @Success 200 {object} []dto.AnalyticsInterval
// @Success 204 ""
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /director/analytics/expert [get]
func (h *DirectorHandler) ExpertAnalytics(w http.ResponseWriter, r *http.Request) {
	expertID := r.URL.Query().Get(expertIDKey)
	if expertID == "" {
		response.BadRequest(w, "Id is empty")
		return
	}

	startTime, err := h.parseTimeQuery(r, startTimeKey)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	endTime, err := h.parseTimeQuery(r, endTimeKey)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	analytics, err := h.directorService.GetExpertAnalytics(expertID, startTime, endTime)
	if err != nil {
		if errors.Is(err, errs.ErrExpertNotExists) {
			response.NotFound(w, "Expert with input ID not found")
			return
		}
		if errors.Is(err, errs.ErrNoRows) {
			response.NoContent(w)
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	analyticsBytes, err := json.Marshal(h.analyticsConverter.MapDomainsToDtos(analytics))
	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, analyticsBytes)
}

// UpdateExpertSkill docs
// @Summary Обновление уровня компетенций у эксперта
// @Security ApiKeyAuth
// @Tags director
// @Description Обновление уровня компетенций у эксперта по его id. Воспользоваться может только директор
// @ID director-expert-skill
// @Produce  json
// @Param input body dto.UpdateExpertSkill true "id эксперта и его новый уровень компетенций"
// @Success 200 {object} response.Body
// @Failure 400,401 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /director/expert_skill [patch]
func (h *DirectorHandler) UpdateExpertSkill(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateExpertSkill
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.directorService.UpdateExpertSkill(input.ExpertID, input.Skill)
	if err != nil {
		response.InternalServerError(w)
		log.Println(err)
		return
	}

	response.OKMessage(w, "Competence skill of expert updated successfully")
}

func (h *DirectorHandler) parseTimeQuery(r *http.Request, key string) (time.Time, error) {
	timeQuery := r.URL.Query().Get(key)
	if timeQuery == "" {
		return time.Time{}, fmt.Errorf("%s is empty", key)
	}

	date, err := time.Parse(time.DateOnly, timeQuery)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
