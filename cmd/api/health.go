package main

import (
	"net/http"
	"time"

	"github.com/ayopedro/seazus-go/lib/logger"
	"github.com/ayopedro/seazus-go/lib/utils"
	"go.uber.org/zap"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	response := utils.APIResponseBody{
		Status:  true,
		Message: "Server is healthy",
		Data: map[string]string{
			"status":      "available",
			"environment": app.config.AppEnv,
			"time":        time.Now().String(),
		},
	}

	err := utils.WriteJSON(w, r, http.StatusOK, response)
	if err != nil {
		logger.Error("json encoding failed", zap.Error(err))
		http.Error(w, "The server encountered a problem", http.StatusInternalServerError)
	}
}
