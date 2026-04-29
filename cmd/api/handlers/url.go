package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	"github.com/ayopedro/seazus-go/cmd/api/response"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
	"go.uber.org/zap"
)

func (h *Handler) GetURLByIdHandler(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	url, err := h.url.GetURL(r.Context(), id, user_id)

	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.WriteError(w, err)
			return
		}
		response.WriteError(w, err)
		return
	}

	result := dto.URL{
		Id:          url.Id,
		Identifier:  url.Identifier,
		Description: url.Description,
		Url:         url.Url,
		ShortUrl:    url.ShortUrl,
		CreatedAt:   dto.JSONTime(url.CreatedAt),
		UpdatedAt:   dto.JSONTime(url.UpdatedAt),
		UserID:      url.UserID,
	}

	response.WriteJSON(w, http.StatusOK, "URL successfully fetched", result)
}

func (h *Handler) GetUserURLSHandler(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	urls, err := h.user.GetUserURLs(r.Context(), user_id)

	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "URLs successfully fetched", urls)
}

func (h *Handler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	payload := &dto.CreateURLPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.WriteError(w, appErrors.ErrInvalidInput)
		return
	}

	if payload.Identifier == "" || payload.Url == "" {
		response.WriteError(w, appErrors.ErrInvalidInput)
		return
	}

	user_id, _ := r.Context().Value(userContextKey).(string)
	short_url, err := h.url.CreateShortURL(
		r.Context(),
		&models.CreateURL{
			Identifier:  payload.Identifier,
			Description: payload.Description,
			Url:         payload.Url,
		},
		user_id,
	)

	if err != nil {
		logger.Error("Error creating short URL", zap.String("err", err.Error()))
		if errors.Is(err, appErrors.ErrConflict) {
			response.WriteError(w, err)
			return
		}
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "URL shortened successfully", short_url)
}
