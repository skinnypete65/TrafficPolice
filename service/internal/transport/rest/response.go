package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseBody struct {
	Message string `json:"message"`
}

func writeResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write(body)
	if err != nil {
		log.Printf("error occurred when write body msg: %v", body)
		return
	}
}

func writeMessage(w http.ResponseWriter, status int, msg string) {
	respBody := ResponseBody{Message: msg}
	body, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error occurred when marshalling response body: %v\n", respBody)
		return
	}

	writeResponse(w, status, body)
}

func badRequest(w http.ResponseWriter, text string) {
	writeMessage(w, http.StatusBadRequest, text)
}

func notFound(w http.ResponseWriter, text string) {
	writeMessage(w, http.StatusNotFound, text)
}

func oKMessage(w http.ResponseWriter, text string) {
	writeMessage(w, http.StatusOK, text)
}

func noContent(w http.ResponseWriter, text string) {
	writeMessage(w, http.StatusNoContent, text)
}

func conflict(w http.ResponseWriter, text string) {
	writeMessage(w, http.StatusConflict, text)
}

func internalServerError(w http.ResponseWriter) {
	writeMessage(w, http.StatusNoContent, "Internal server error")
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeMessage(w, http.StatusUnauthorized, "Unauthorized")
}
