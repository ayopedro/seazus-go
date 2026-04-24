package handler

import (
	"net/http"

	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/utils"
	"go.uber.org/zap"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	response := utils.APIResponseBody{
		Status:  true,
		Message: "Welcome to the Seazus(Go) API",
		Data:    []string{},
	}

	err := utils.WriteJSON(w, r, http.StatusOK, response)
	if err != nil {
		logger.Error("json encoding failed", zap.Error(err))
		utils.WriteError(w, r, http.StatusInternalServerError, err)
	}
}
