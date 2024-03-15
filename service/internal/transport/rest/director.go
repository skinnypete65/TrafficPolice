package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	expertIDKey  = "id"
	startTimeKey = "start_time"
	endTimeKey   = "end_time"
)

type DirectorHandler struct {
	directorService    services.DirectorService
	caseConverter      *converter.CaseConverter
	analyticsConverter *converter.AnalyticsConverter
}

func NewDirectorHandler(
	directorService services.DirectorService,
	caseConverter *converter.CaseConverter,
	analyticsConverter *converter.AnalyticsConverter,
) *DirectorHandler {
	return &DirectorHandler{
		directorService:    directorService,
		caseConverter:      caseConverter,
		analyticsConverter: analyticsConverter,
	}
}

// GetCases docs
// @Summary Получение состояния случаев
// @Security ApiKeyAuth
// @Tags director
// @Description Получение состояния случаев. Воспользоваться может только директор
// @ID director-cases-get
// @Produce  json
// @Success 200 {object} []dto.CaseStatus
// @Success 204 ""
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /director/cases [get]
func (h *DirectorHandler) GetCases(w http.ResponseWriter, r *http.Request) {
	cases, err := h.directorService.GetCases()
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			response.NoContent(w)
			return
		}
		log.Println(err)
		response.InternalServerError(w)
	}

	casesBytes, err := json.Marshal(h.caseConverter.MapCaseStatusesToDto(cases))
	if err != nil {
		response.InternalServerError(w)
	}

	response.WriteResponse(w, http.StatusOK, casesBytes)
}

// ExpertAnalytics docs
// @Summary Получение аналитики проверяющих специалистов по промежуткам времени
// @Security ApiKeyAuth
// @Tags director
// @Description Получить количество всех случаев, правильно решенных случае, неправильно решенных случаев,
// неизвестных случаев и максимальное количество подряд решенных задач. Воспользоваться может только директор
// @ID director-analytics-expert
// @Produce  json
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
