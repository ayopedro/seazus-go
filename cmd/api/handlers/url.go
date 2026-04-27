package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
	"go.uber.org/zap"
)

func (h *handler) GetURLByIdHandler(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	url, err := h.urlService.GetURL(r.Context(), id, user_id)

	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			utils.WriteError(w, r, err)
		}
		utils.WriteError(w, r, err)
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "URL successfully fetched",
		Data:    &url,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func (h *handler) GetUserURLSHandler(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	urls, err := h.userService.GetUserURLs(r.Context(), user_id)

	if err != nil {
		utils.WriteError(w, r, err)
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "URLs successfully fetched",
		Data:    urls,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func (h *handler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.CreateURLPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, appErrors.ErrInvalidInput)
		return
	}

	if payload.Identifier == "" || payload.Url == "" {
		utils.WriteError(w, r, appErrors.ErrInvalidInput)
		return
	}

	user_id, _ := r.Context().Value(userContextKey).(string)
	short_url, err := h.urlService.CreateShortURL(r.Context(), payload, user_id)

	if err != nil {
		logger.Error("Error creating short URL", zap.String("err", err.Error()))
		if errors.Is(err, appErrors.ErrConflict) {
			utils.WriteError(w, r, err)
			return
		}
		utils.WriteError(w, r, err)
		return
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "URL shortened successfully",
		Data:    short_url,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}
