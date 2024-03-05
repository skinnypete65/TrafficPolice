package transport

import (
	"TrafficPolice/internal/models"
	"TrafficPolice/internal/services"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	casesDir            = "cases"
	caseContentImageKey = "image"
	caseIDPathValue     = "id"
)

var mapping = map[string]func(c *models.Case, value any) error{
	"transport_chars":   setTransportChars,
	"transport_numbers": setTransportNums,
	"transport_region":  setTransportRegion,
	"camera_id":         setCameraID,
	"violation_id":      setViolationID,
	"violation_value":   setViolationValue,
	"skill_value":       setSkillValue,
	"datetime":          setDatetime,
}

func setTransportChars(c *models.Case, value any) error {
	c.Transport.Chars = value.(string)
	return nil
}

func setTransportNums(c *models.Case, value any) error {
	c.Transport.Num = value.(string)
	return nil
}

func setTransportRegion(c *models.Case, value any) error {
	c.Transport.Region = value.(string)
	return nil
}

func setCameraID(c *models.Case, value any) error {
	c.Camera.ID = value.(string)
	return nil
}

func setViolationID(c *models.Case, value any) error {
	c.Violation.ID = value.(string)
	return nil
}

func setViolationValue(c *models.Case, value any) error {
	c.ViolationValue = value.(string)
	return nil
}

func setSkillValue(c *models.Case, value any) error {
	c.RequiredSkill = value.(int)
	return nil
}

func setDatetime(c *models.Case, value any) error {
	t, err := time.Parse(time.RFC3339, value.(string))
	if err != nil {
		return err
	}
	c.Date = t
	return nil
}

type CaseHandler struct {
	service    services.CaseService
	imgService services.ImgService
}

func NewCaseHandler(service services.CaseService, imgService services.ImgService) *CaseHandler {

	return &CaseHandler{
		service:    service,
		imgService: imgService,
	}
}

func (h *CaseHandler) AddCase(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(buf)
	log.Println()

	inputCase, err := parseCase(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.AddCase(inputCase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("Case added successfully"))
	if err != nil {
		log.Println(err)
	}
}

func parseCase(payload []byte) (*models.Case, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("payload is empty")
	}
	transportCase := &models.Case{}
	payload = payload[2:]

	for len(payload) > 0 {
		keySize := binary.BigEndian.Uint16(payload[:2])
		payload = payload[2:]

		valueSize := binary.BigEndian.Uint16(payload[:2])
		payload = payload[2:]

		valueType := payload[0]
		payload = payload[1:]

		keyValue := payload[:keySize]
		value := payload[keySize : keySize+valueSize]
		payload = payload[keySize+valueSize:]

		f := mapping[string(keyValue)]

		var err error
		if valueType == 0 {
			err = f(transportCase, string(value))
		}

		if err != nil {
			return nil, err
		}
	}

	return transportCase, nil
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
