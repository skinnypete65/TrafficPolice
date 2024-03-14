package rest

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport/rabbitmq"
	"TrafficPolice/internal/transport/rest/dto"
	"TrafficPolice/internal/transport/rest/middlewares"
	"TrafficPolice/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	expertsDir            = "experts"
	expertContentImageKey = "image"
	expertIDPathValue     = "id"
)

type ExpertHandler struct {
	imgService            services.ImgService
	expertService         services.ExpertService
	ratingService         services.RatingService
	finePublisher         *rabbitmq.FinePublisher
	caseConverter         *converter.CaseConverter
	caseDecisionConverter *converter.CaseDecisionConverter
}

func NewExpertHandler(
	imgService services.ImgService,
	expertService services.ExpertService,
	ratingService services.RatingService,
	finePublisher *rabbitmq.FinePublisher,
	caseConverter *converter.CaseConverter,
	caseDecisionConverter *converter.CaseDecisionConverter,
) *ExpertHandler {
	return &ExpertHandler{
		imgService:            imgService,
		expertService:         expertService,
		ratingService:         ratingService,
		finePublisher:         finePublisher,
		caseConverter:         caseConverter,
		caseDecisionConverter: caseDecisionConverter,
	}
}

// UploadExpertImg docs
// @Summary Добавление фотографии к профилю эксперта
// @Security ApiKeyAuth
// @Tags expert
// @Description Добавление фотографии к профилю эксперта. Может воспользоваться директор или эксперт
// @ID expert-image-upload
// @Accept   multipart/form-data
// @Produce  json
// @Param id query string true "id эксперта"
// @Param file formData file true "Фотография эксперта"
// @Success 200 {object} response.Body
// @Failure 400,401 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /expert/{id}/img [post]
func (h *ExpertHandler) UploadExpertImg(w http.ResponseWriter, r *http.Request) {
	expertID := r.PathValue(expertIDPathValue)
	if expertID == "" {
		response.BadRequest(w, "Bad expert id")
		return
	}

	file, header, err := parseMultipartForm(r, expertContentImageKey)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	contentType := header.Header.Get(contentTypeKey)
	extension, err := getImgExtension(contentType)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	imgFilePath := fmt.Sprintf("%s/%s.%s", expertsDir, expertID, extension)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error while reading fileBytes: %v\n", fileBytes)
		response.InternalServerError(w)
		return
	}

	err = h.imgService.SaveImg(fileBytes, imgFilePath)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.OKMessage(w, "Successfully uploaded image")
}

// GetExpertImg docs
// @Summary Получение фотографии эксперта
// @Security ApiKeyAuth
// @Tags expert
// @Description Получение фотографии эксперта по его id. Воспользоваться могут эксперт или директор
// @ID expert-image-get
// @Accept   multipart/form-data
// @Produce  json
// @Param id query string true "id эксперта"
// @Success 200 {file} formData
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /expert/{id}/img [get]
func (h *ExpertHandler) GetExpertImg(w http.ResponseWriter, r *http.Request) {
	expertID := r.PathValue(expertIDPathValue)
	if expertID == "" {
		response.BadRequest(w, "Bad expert id")
		return
	}

	file, err := h.imgService.GetImgFilePath(expertsDir, expertID)
	if err != nil {
		if errors.Is(err, errs.ErrNoImage) {
			response.NotFound(w, "Image with input expert id not found")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	http.ServeFile(w, r, file)
}

// GetCaseForExpert docs
// @Summary Получение случая для оценки экспертом
// @Security ApiKeyAuth
// @Tags expert
// @Description Получение случая для оценки экспертом. Воспользоваться могут эксперт или директор
// @ID expert-get-case
// @Produce  json
// @Success 200 {file} dto.Case
// @Success 204 ""
// @Failure 401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /expert/get_case [get]
func (h *ExpertHandler) GetCaseForExpert(w http.ResponseWriter, r *http.Request) {
	tokenInfo := r.Context().Value(middlewares.TokenInfoKey).(tokens.TokenInfo)

	c, err := h.expertService.GetCase(tokenInfo.UserID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotExists) {
			response.NotFound(w, "Expert not found")
		}
		if errors.Is(err, errs.ErrNoNotSolvedCase) || errors.Is(err, errs.ErrNoCase) {
			response.NoContent(w)
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	cBytes, err := json.Marshal(h.caseConverter.MapDomainToDto(c))
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, cBytes)
}

// SetCaseDecision docs
// @Summary Оценка случая экспертом
// @Security ApiKeyAuth
// @Tags expert
// @Description Установка оценки случая экспертом. Воспользоваться может только эксперт
// @ID expert-set-decision
// @Accept   json
// @Produce  json
// @Param input body dto.Decision true "id случая и решение эксперта"
// @Success 200 {file} response.Body
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /expert/decision [post]
func (h *ExpertHandler) SetCaseDecision(w http.ResponseWriter, r *http.Request) {
	tokenInfo := r.Context().Value(middlewares.TokenInfoKey).(tokens.TokenInfo)

	var decision dto.Decision
	err := json.NewDecoder(r.Body).Decode(&decision)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	expert, err := h.expertService.GetExpertByUserID(tokenInfo.UserID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotExists) {
			response.NotFound(w, "Input expert not found")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	caseDecision, err := h.expertService.SetCaseDecision(
		h.caseDecisionConverter.MapDtoToDomain(decision, expert),
	)

	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	if caseDecision.IsSolved {
		err = h.ratingService.SetRating(caseDecision)
		if err != nil {
			log.Println(err)
			response.InternalServerError(w)
		}
	}

	if caseDecision.ShouldSendFine {
		caseInfo, err := h.expertService.GetCaseWithPersonInfo(decision.CaseID)
		if err != nil {
			log.Println(err)
			return
		}
		c := h.caseConverter.MapCaseWithPersonToDTO(caseInfo)
		image, extension, err := h.getCaseImage(caseInfo.ID)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.finePublisher.PublishFineNotification(dto.CaseWithImage{
			Case:           c,
			Image:          image,
			ImageExtension: extension,
		})
		if err != nil {
			log.Println(err)
			return
		}
	}

	response.OKMessage(w, "Decision accepted")
}

func (h *ExpertHandler) getCaseImage(caseID string) ([]byte, string, error) {
	file, err := h.imgService.GetImgFilePath(casesDir, caseID)
	if err != nil {
		return nil, "", err
	}

	img, err := os.ReadFile(file)
	if err != nil {
		return nil, "", err
	}

	dotIdx := strings.LastIndex(file, ".")
	return img, file[dotIdx+1:], nil
}
