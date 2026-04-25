package handlers

import (
	"errors"
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/models"
)

func (h *handler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(userContextKey).(*models.User)

	response := types.APIResponseBody{
		Status:  true,
		Message: "User profile successfully fetched",
		Data:    &user,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func (h *handler) GetURLById(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	url, err := h.urlService.GetURL(r.Context(), id, user_id)

	if err != nil {
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			utils.WriteError(w, r, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, r, http.StatusInternalServerError, err)
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "URL successfully fetched",
		Data:    &url,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
