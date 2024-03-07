package transport

import (
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/internal/transport/dto"
	"TrafficPolice/internal/transport/middlewares"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

const (
	expertsDir            = "experts"
	expertContentImageKey = "image"
	expertIDPathValue     = "id"
)

type ExpertHandler struct {
	imgService    services.ImgService
	expertService services.ExpertService
}

func NewExpertHandler(imgService services.ImgService, expertService services.ExpertService) *ExpertHandler {
	return &ExpertHandler{
		imgService:    imgService,
		expertService: expertService,
	}
}

func (h *ExpertHandler) UploadExpertImg(w http.ResponseWriter, r *http.Request) {
	expertID := r.PathValue(expertIDPathValue)
	if expertID == "" {
		http.Error(w, "id is empty", http.StatusBadRequest)
		return
	}

	file, header, err := parseMultipartForm(r, expertContentImageKey)
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

	imgFilePath := fmt.Sprintf("%s/%s.%s", expertsDir, expertID, extension)
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

func (h *ExpertHandler) GetExpertImg(w http.ResponseWriter, r *http.Request) {
	expertID := r.PathValue(expertIDPathValue)
	if expertID == "" {
		http.Error(w, "bad expert id", http.StatusBadRequest)
		return
	}

	pattern := fmt.Sprintf("%s/%s.*", expertsDir, expertID)
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.ServeFile(w, r, files[0])
}

func (h *ExpertHandler) GetCaseForExpert(w http.ResponseWriter, r *http.Request) {
	tokenInfo := r.Context().Value(middlewares.TokenInfoKey).(tokens.TokenInfo)
	log.Println(tokenInfo)

	c, err := h.expertService.GetCase(tokenInfo.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cDto := dto.Case{
		ID: c.ID,
		Transport: dto.Transport{
			ID:     c.Transport.ID,
			Chars:  c.Transport.Chars,
			Num:    c.Transport.Num,
			Region: c.Transport.Region,
			Person: dto.Person{
				ID: c.Transport.Person.ID,
			},
		},
		Camera: dto.Camera{
			ID:           c.Camera.ID,
			CameraTypeID: c.Camera.CameraTypeID,
			Latitude:     c.Camera.Latitude,
			Longitude:    c.Camera.Longitude,
			ShortDesc:    c.Camera.ShortDesc,
		},
		Violation: dto.Violation{
			ID:         c.Violation.ID,
			Name:       c.Violation.Name,
			FineAmount: c.Violation.FineAmount,
		},
		ViolationValue: c.ViolationValue,
		RequiredSkill:  c.RequiredSkill,
		IsSolved:       c.IsSolved,
		FineDecision:   c.FineDecision,
	}

	cBytes, err := json.Marshal(cDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(cBytes)
	if err != nil {
		log.Println()
	}
}
