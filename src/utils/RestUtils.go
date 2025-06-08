package utils

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	RespondJSON(w, map[string]string{"error": message}, status)
}

func RespondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(data)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		panic(err)
	}
}
