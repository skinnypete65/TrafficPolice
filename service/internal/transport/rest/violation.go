package rest

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/response"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
)

const (
	violationSheet   = "Лист1"
	violationFileKey = "file"
)

type ViolationHandler struct {
	service services.ViolationService
}

func NewViolationHandler(service services.ViolationService) *ViolationHandler {
	return &ViolationHandler{service: service}
}

// InsertViolations docs
// @Summary Ввод информации о правонарушениях
// @Security ApiKeyAuth
// @Tags violation
// @Description Принимает excel файл в формате по столбикам: Название правонарушения:Размер штрафа. Только директор может загрузить файл
// @ID insert-violations
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "Excel файл с информацией о правонарушениях"
// @Success 200 {object} response.Body
// @Failure 400,401 {object} response.Body
// @Failure 500 {object} response.Body
// @Failure default {object} response.Body
// @Router /violations [post]
func (h *ViolationHandler) InsertViolations(w http.ResponseWriter, r *http.Request) {
	maxMemory := int64(10 << 30)

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		text := fmt.Sprintf("Error while ParseMultipartForm: %v", err)
		response.BadRequest(w, text)
		return
	}

	// retrieve file from posted form-data
	formFile, _, err := r.FormFile(violationFileKey)
	if err != nil {
		text := fmt.Sprintf("Error retrieving file from form-data: %v\n", err)
		response.BadRequest(w, text)
		return
	}
	defer formFile.Close()

	f, err := excelize.OpenReader(formFile)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	rows, err := f.GetRows(violationSheet)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	violations := h.parseViolations(rows)

	err = h.service.InsertViolations(violations)
	if err != nil {
		response.InternalServerError(w)
		return
	}
}

func (h *ViolationHandler) parseViolations(rows [][]string) []*domain.Violation {
	violations := make([]*domain.Violation, 0)

	for _, row := range rows {
		isValid := true
		for _, val := range row {
			if val == "" {
				isValid = false
			}
			break
		}
		if !isValid {
			continue
		}
		fineAmount, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}

		v := &domain.Violation{
			Name:       row[0],
			FineAmount: fineAmount,
		}

		violations = append(violations, v)
	}

	return violations
}
