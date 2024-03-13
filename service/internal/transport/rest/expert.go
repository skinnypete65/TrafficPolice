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
	finePublisher         *rabbitmq.FinePublisher
	caseConverter         *converter.CaseConverter
	caseDecisionConverter *converter.CaseDecisionConverter
}

func NewExpertHandler(
	imgService services.ImgService,
	expertService services.ExpertService,
	finePublisher *rabbitmq.FinePublisher,
	caseConverter *converter.CaseConverter,
	caseDecisionConverter *converter.CaseDecisionConverter,
) *ExpertHandler {
	return &ExpertHandler{
		imgService:            imgService,
		expertService:         expertService,
		finePublisher:         finePublisher,
		caseConverter:         caseConverter,
		caseDecisionConverter: caseDecisionConverter,
	}
}

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

	shouldSendFine, err := h.expertService.SetCaseDecision(
		h.caseDecisionConverter.MapDtoToDomain(decision, expert),
	)

	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	if shouldSendFine {
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
