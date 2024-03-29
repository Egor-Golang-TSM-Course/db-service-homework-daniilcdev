package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg interface{}) {
	if code > 499 {
		log.Println("Responsing with 5XX error:", msg)
	}

	RespondWithJSON(w, code, errResponse{Error: fmt.Sprintf("%v", msg)})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v\n", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	SetHeaders(w, code)
	if _, err = w.Write(data); err != nil {
		log.Printf("failed to write response to buffer: code=%d\n", code)
	}
}

func SetHeaders(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
}
