package rest

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/transport/rest/response"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
)

const (
	contactInfoSheet   = "Лист1"
	contactInfoFileKey = "file"
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
		response.BadRequest(w, err.Error())
		return
	}

	// retrieve file from posted form-data
	formFile, _, err := r.FormFile(contactInfoFileKey)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	defer formFile.Close()

	f, err := excelize.OpenReader(formFile)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	rows, err := f.GetRows(contactInfoSheet)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
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
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.OKMessage(w, "Contact info added successfully")
}
