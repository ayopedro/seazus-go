package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ayopedro/seazus-go/internal/logger"
	"go.uber.org/zap"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, response APIResponseBody) error {
	js, err := json.Marshal(response)
	if err != nil {
		logger.Error("json encoding failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) error {
	response := APIResponseBody{
		Status:  false,
		Message: err.Error(),
	}

	js, err := json.Marshal(response)
	if err != nil {
		logger.Error("json encoding failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
