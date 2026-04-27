package handlers

import (
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/logger"
	"go.uber.org/zap"
)

func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  true,
		Message: "Welcome to the Seazus(Go) API",
		Data:    []string{},
	}

	err := utils.WriteJSON(w, r, http.StatusOK, response)
	if err != nil {
		logger.Error("json encoding failed", zap.Error(err))
		utils.WriteError(w, r, err)
	}
}
