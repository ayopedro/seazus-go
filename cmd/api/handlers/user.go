package handlers

import (
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/models"
)

func (h *handler) GetMyProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(userContextKey).(*models.User)

	response := types.APIResponseBody{
		Status:  true,
		Message: "User profile successfully fetched",
		Data:    &user,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
