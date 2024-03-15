package rest

import (
	"TrafficPolice/internal/camera"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/transport/rest/response"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

const (
	casesDir            = "cases"
	caseContentImageKey = "image"
	caseIDPathValue     = "id"
)

type CaseHandler struct {
	caseService   service.CaseService
	imgService    service.ImgService
	cameraService service.CameraService
	caseConverter *converter.CaseConverter
	cameraParser  *camera.Parser
}

func NewCaseHandler(
	service service.CaseService,
	imgService service.ImgService,
	cameraService service.CameraService,
	caseConverter *converter.CaseConverter,
	cameraParser *camera.Parser,
) *CaseHandler {

	return &CaseHandler{
		caseService:   service,
		imgService:    imgService,
		cameraService: cameraService,
		caseConverter: caseConverter,
		cameraParser:  cameraParser,
	}
}

// AddCase docs
// @Summary Добавление информации о проишествии
// @Security ApiKeyAuth
// @Tags case
// @Description Принимает бинарную строку в описанном формате. Добавить проишествие может только камера
// @ID case-add
// @Accept   application/octet-stream
// @Produce  json
// @Success 200 {object} response.IDResponse
// @Failure 400,401 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /case [post]
func (h *CaseHandler) AddCase(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	inputCase, err := h.cameraParser.ParseCameraInfo(buf)
	if err != nil {
		if errors.Is(err, errs.ErrEmptyPayload) {
			response.BadRequest(w, "Binary string is empty")
			return
		}
		if errors.Is(err, errs.ErrUnknownCameraID) {
			response.BadRequest(w, "Cannot parse camera id")
			return
		}
		if errors.Is(err, errs.ErrUnknownCameraType) {
			response.BadRequest(w, "Unknown camera type")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	caseID, err := h.caseService.AddCase(h.caseConverter.MapDtoToDomain(inputCase))
	if err != nil {
		if errors.Is(err, errs.ErrNoTransport) {
			response.BadRequest(w, "No transport by input credentials")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.IdResponse(w, caseID)
}

// UploadCaseImg docs
// @Summary Добавление фотографии к проишествию
// @Security ApiKeyAuth
// @Tags case
// @Description Принимает фотографию и сохраняет ее по переданному id. Добавить фотографию может только камера
// @ID case-image-upload
// @Accept   multipart/form-data
// @Produce  json
// @Param id query string true "id камеры"
// @Param file formData file true "Фотография проишествия"
// @Success 200 {object} response.Body
// @Failure 400,401 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /case/{id}/img [post]
func (h *CaseHandler) UploadCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue(caseIDPathValue)
	if caseID == "" {
		response.BadRequest(w, "id is empty")
		return
	}

	file, header, err := parseMultipartForm(r, caseContentImageKey)
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

	imgFilePath := fmt.Sprintf("%s/%s.%s", casesDir, caseID, extension)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error while reading fileBytes: %v\n", fileBytes)
		response.InternalServerError(w)
		return
	}

	err = h.imgService.SaveImg(fileBytes, imgFilePath)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.OKMessage(w, "Successfully uploaded image")
}

// GetCaseImg docs
// @Summary Получение фотографии проишествия
// @Security ApiKeyAuth
// @Tags case
// @Description Получение фотографии проишествия по id прошествия. Воспользоваться могут эксперт или директор
// @ID case-image-get
// @Accept   multipart/form-data
// @Produce  json
// @Param id query string true "id камеры"
// @Success 200 {file} formData
// @Failure 400,401,404 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /case/{id}/img [get]
func (h *CaseHandler) GetCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue(caseIDPathValue)
	if caseID == "" {
		response.BadRequest(w, "bad case id")
		return
	}

	file, err := h.imgService.GetImgFilePath(casesDir, caseID)
	if err != nil {
		if errors.Is(err, errs.ErrNoImage) {
			response.NotFound(w, "Image with input case id not found")
			return
		}
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	http.ServeFile(w, r, file)
}
