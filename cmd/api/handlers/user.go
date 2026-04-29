package handlers

import (
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	"github.com/ayopedro/seazus-go/cmd/api/response"
)

func (h *Handler) GetMyProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userContextKey).(string)

	user, err := h.user.GetUserProfile(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err)
		return
	}

	result := dto.User{
		Id:         user.Id,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		IsVerified: user.IsVerified,
		CreatedAt:  dto.JSONTime(user.CreatedAt),
		UpdatedAt:  dto.JSONTime(user.UpdatedAt),
	}

	response.WriteJSON(w, http.StatusOK, "User profile successfully fetched", result)
}
