package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, response APIResponseBody) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponseBody{
		Status:  false,
		Message: err.Error(),
	}

	return json.NewEncoder(w).Encode(response)
}
