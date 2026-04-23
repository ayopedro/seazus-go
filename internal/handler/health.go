package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/utils"
	"go.uber.org/zap"
)

type Handler struct {
	AppConfig *config.Config
	DB        *sql.DB
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := utils.APIResponseBody{
		Status:  true,
		Message: "Server is healthy",
		Data: map[string]string{
			"status":      "available",
			"environment": h.AppConfig.AppEnv,
			"time":        time.Now().String(),
		},
	}

	err := utils.WriteJSON(w, r, http.StatusOK, response)
	if err != nil {
		logger.Error("json encoding failed", zap.Error(err))
		http.Error(w, "The server encountered a problem", http.StatusInternalServerError)
	}
}
