package rest

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/camera"
	"TrafficPolice/internal/converter"
	"TrafficPolice/internal/services"
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
	caseService   services.CaseService
	imgService    services.ImgService
	cameraService services.CameraService
	caseConverter *converter.CaseConverter
	cameraParser  *camera.Parser
}

func NewCaseHandler(
	service services.CaseService,
	imgService services.ImgService,
	cameraService services.CameraService,
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
