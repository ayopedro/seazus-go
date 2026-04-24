package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/types"
	"github.com/ayopedro/seazus-go/internal/utils"
)

type Handler struct {
	AppConfig *config.Config
	DB        *sql.DB
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  true,
		Message: "Server is healthy",
		Data: map[string]string{
			"status":      "available",
			"environment": h.AppConfig.AppEnv,
			"time":        time.Now().String(),
		},
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
