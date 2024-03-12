package rest

import (
	"TrafficPolice/internal/camera"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/dto"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
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
	cameraParser  *camera.Parser
}

func NewCaseHandler(
	service services.CaseService,
	imgService services.ImgService,
	cameraService services.CameraService,
) *CaseHandler {

	return &CaseHandler{
		caseService:   service,
		imgService:    imgService,
		cameraService: cameraService,
		cameraParser:  camera.NewParser(cameraService),
	}
}

func (h *CaseHandler) AddCase(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("BUF:")
	log.Println(buf)

	inputCase, err := h.cameraParser.ParseCameraInfo(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.caseService.AddCase(mapCaseDTOtoDomain(inputCase))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CaseHandler) UploadCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue(caseIDPathValue)
	if caseID == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
		return
	}

	file, header, err := parseMultipartForm(r, caseContentImageKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get(contentTypeKey)
	extension, err := getImgExtension(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgFilePath := fmt.Sprintf("%s/%s.%s", casesDir, caseID, extension)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error while reading fileBytes: %v\n", fileBytes)
		return
	}

	err = h.imgService.SaveImg(fileBytes, imgFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Succesfully uploaded image")
}

func (h *CaseHandler) GetCaseImg(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue(caseIDPathValue)
	if caseID == "" {
		http.Error(w, "bad case id", http.StatusBadRequest)
		return
	}

	pattern := fmt.Sprintf("%s/%s.*", casesDir, caseID)
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.ServeFile(w, r, files[0])
}

func mapCaseDTOtoDomain(c dto.Case) domain.Case {
	return domain.Case{
		Transport: domain.Transport{
			Chars:  c.Transport.Chars,
			Num:    c.Transport.Num,
			Region: c.Transport.Region,
		},
		Camera: domain.Camera{
			ID: c.Camera.ID,
		},
		Violation: domain.Violation{
			ID: c.Violation.ID,
		},
		ViolationValue: c.ViolationValue,
		RequiredSkill:  c.RequiredSkill,
		Date:           c.Date,
	}
}
