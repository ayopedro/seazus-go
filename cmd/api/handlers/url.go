package handlers

import (
	"errors"
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/common/types"
)

func (h *handler) GetURLById(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	url, err := h.urlService.GetURL(r.Context(), id, user_id)

	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			utils.WriteError(w, r, http.StatusNotFound, err)
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

func (h *handler) GetUserURLS(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	urls, err := h.userService.GetUserURLs(r.Context(), user_id)

	if err != nil {
		utils.WriteError(w, r, http.StatusInternalServerError, err)
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "URLs successfully fetched",
		Data:    urls,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
