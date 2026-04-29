package handlers

import (
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/response"
)

func (h *Handler) GetMyProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userContextKey).(string)

	user, err := h.user.GetUserProfile(r.Context(), userID)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	result := response.APIResponseBody{
		Status:  true,
		Message: "User profile successfully fetched",
		Data:    &user,
	}

	response.WriteJSON(w, r, http.StatusOK, result)
}
