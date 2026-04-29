package response

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ResponseBody{
		Status:  statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
