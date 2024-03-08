package transport

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
)

const (
	violationSheet = "Лист1"
)

type ViolationHandler struct {
	service services.ViolationService
}

func NewViolationHandler(service services.ViolationService) *ViolationHandler {
	return &ViolationHandler{service: service}
}

func (h *ViolationHandler) InsertViolations(w http.ResponseWriter, r *http.Request) {
	maxMemory := int64(10 << 30)

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		log.Printf("Error while ParseMultipartForm: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// retrieve file from posted form-data
	formFile, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from form-data: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	f, err := excelize.OpenReader(formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := f.GetRows(violationSheet)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	err = h.service.InsertViolations(violations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
