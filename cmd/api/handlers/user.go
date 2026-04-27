package handlers

import (
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	"github.com/ayopedro/seazus-go/internal/common/types"
)

func (h *handler) GetMyProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userContextKey).(string)

	user, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, r, err)
		return
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "User profile successfully fetched",
		Data:    &user,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
