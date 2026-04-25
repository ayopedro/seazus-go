package handlers

import (
	"net/http"
	"time"

	utils "github.com/ayopedro/seazus-go/internal/common"
	"github.com/ayopedro/seazus-go/internal/common/types"
)

func (h *handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  true,
		Message: "Server is healthy",
		Data: map[string]string{
			"status":      "available",
			"environment": h.config.AppEnv,
			"time":        time.Now().String(),
		},
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
