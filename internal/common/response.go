package common

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/logger"
	"go.uber.org/zap"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, response types.APIResponseBody) error {
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

func WriteError(w http.ResponseWriter, r *http.Request, err error) error {
	status := appErrors.StatusCode(err)

	response := types.APIResponseBody{
		Status:  false,
		Message: err.Error(),
	}

	if errors.Is(err, appErrors.ErrInternal) {
		response.Message = http.StatusText(http.StatusInternalServerError)
	}

	js, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		logger.Error("json encoding failed", zap.Error(marshalErr))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return marshalErr
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
