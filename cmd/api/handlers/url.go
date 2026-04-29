package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	"github.com/ayopedro/seazus-go/cmd/api/response"
	"github.com/ayopedro/seazus-go/internal/apperrors"
	"github.com/ayopedro/seazus-go/internal/logger"
	"github.com/ayopedro/seazus-go/internal/models"
	"go.uber.org/zap"
)

func (h *Handler) GetURLByIdHandler(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	url, err := h.url.GetURL(r.Context(), id, user_id)

	if err != nil {
		if errors.Is(err, apperrors.ErrURLNotFound) {
			response.WriteError(w, http.StatusNotFound, err)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
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
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var result []dto.URL
	for _, url := range urls {
		u := dto.URL{
			Id:          url.Id,
			Identifier:  url.Identifier,
			Url:         url.Url,
			ShortUrl:    url.ShortUrl,
			Description: url.Description,
			CreatedAt:   dto.JSONTime(url.CreatedAt),
			UpdatedAt:   dto.JSONTime(url.UpdatedAt),
			UserID:      url.UserID,
		}

		result = append(result, u)
	}

	response.WriteJSON(w, http.StatusOK, "URLs successfully fetched", result)
}

func (h *Handler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	payload := &dto.CreateURLPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidPayload)
		return
	}

	if payload.Identifier == "" || payload.Url == "" {
		response.WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidPayload)
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
		if errors.Is(err, apperrors.ErrIdentifierTaken) {
			response.WriteError(w, http.StatusConflict, err)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "URL shortened successfully", short_url)
}

func (h *Handler) UpdateURLHandler(w http.ResponseWriter, r *http.Request) {
	payload := &dto.UpdateURLPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.WriteError(w, http.StatusBadRequest, apperrors.ErrInvalidPayload)
		return
	}

	uID, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	updatedURL, err := h.url.UpdateURL(
		r.Context(),
		id,
		&models.UpdateURL{
			Identifier:  payload.Identifier,
			Description: payload.Description,
			Url:         payload.Url,
		},
		uID,
	)

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "URL updated successfully", updatedURL)
}

func (h *Handler) DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	uID, _ := r.Context().Value(userContextKey).(string)
	id := r.PathValue("id")
	err := h.url.DeleteURL(
		r.Context(),
		id,
		uID,
	)

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "URL deleted successfully", nil)
}
