package rest

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
)

const (
	contactInfoSheet = "Лист1"
)

type ContactInfoHandler struct {
	service services.ContactInfoService
}

func NewContactInfoHandler(service services.ContactInfoService) *ContactInfoHandler {
	return &ContactInfoHandler{service: service}
}

func (h *ContactInfoHandler) InsertContactInfo(w http.ResponseWriter, r *http.Request) {
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

	rows, err := f.GetRows(contactInfoSheet)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := make(map[string][]*domain.Transport)

	for _, row := range rows {
		transport := &domain.Transport{
			Chars:  row[0],
			Num:    row[1],
			Region: row[2],
		}

		person := &domain.Person{
			PhoneNum: row[3],
			Email:    row[4],
			VkID:     row[5],
			TgID:     row[6],
		}
		transport.Person = person

		if _, ok := m[person.PhoneNum]; !ok {
			m[person.PhoneNum] = make([]*domain.Transport, 0)
		}
		m[person.PhoneNum] = append(m[person.PhoneNum], transport)
	}

	err = h.service.InsertContactInfo(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte("Added successfully"))
	if err != nil {
		log.Println(err)
	}
}
