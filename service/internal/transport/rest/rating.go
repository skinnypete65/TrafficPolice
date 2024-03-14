package rest

import (
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RatingHandler struct {
	ratingService   services.RatingService
	ratingConverter converter.RatingConverter
}

func NewRatingHandler(
	ratingService services.RatingService,
	ratingConverter *converter.RatingConverter,
) *RatingHandler {
	return &RatingHandler{
		ratingService:   ratingService,
		ratingConverter: *ratingConverter,
	}
}

// GetRating docs
// @Summary Получение рейтинга экспертов
// @Security ApiKeyAuth
// @Tags rating
// @Description Получение рейтинга экспертов. Воспользоваться могут эксперт или директор
// @ID rating-get
// @Produce  json
// @Success 200 {object} []dto.RatingInfo
// @Success 204 ""
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /rating [get]
func (h *RatingHandler) GetRating(w http.ResponseWriter, r *http.Request) {
	rating, err := h.ratingService.GetRating()
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			response.NoContent(w)
			return
		}
		log.Println(err)
		response.InternalServerError(w)
	}

	ratingBytes, err := json.Marshal(h.ratingConverter.MapSliceDomainToDto(rating))
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
	}

	response.WriteResponse(w, http.StatusOK, ratingBytes)
}
