package rest

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type DirectorHandler struct {
	directorService services.DirectorService
	caseConverter   *converter.CaseConverter
}

func NewDirectorHandler(
	directorService services.DirectorService,
	caseConverter *converter.CaseConverter,
) *DirectorHandler {
	return &DirectorHandler{
		directorService: directorService,
		caseConverter:   caseConverter,
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
