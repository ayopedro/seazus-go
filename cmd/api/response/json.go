package response

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, message string, data any) error {
	response := ResponseBody{
		Status:  status >= 200 && status < 300,
		Message: message,
		Data:    data,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	return err
}
