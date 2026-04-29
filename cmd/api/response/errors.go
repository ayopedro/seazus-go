package response

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/logger"
	"go.uber.org/zap"
)

func WriteError(w http.ResponseWriter, err error) error {
	status := appErrors.StatusCode(err)

	response := ResponseBody{
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
