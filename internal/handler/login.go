package handler

import (
	"net/http"

	"github.com/ayopedro/seazus-go/internal/types"
	"github.com/ayopedro/seazus-go/internal/utils"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  true,
		Message: "You are on the login route",
	}

	utils.WriteJSON(w, r, http.StatusCreated, response)
}
