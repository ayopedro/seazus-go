package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ayopedro/seazus-go/cmd/api/response"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, "Welcome to the Seazus(Go) API", nil)
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{
		"status": "available",
		"time":   time.Now().String(),
	}

	response.WriteJSON(w, http.StatusOK, "", result)
}

func (h *Handler) ShortURLRedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("short_url")
	original_url, err := h.url.GetOriginalURL(r.Context(), path)

	if err != nil {
		response.WriteError(w, err)
		return
	}

	http.Redirect(w, r, original_url, http.StatusSeeOther)
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	response.WriteError(w, errors.New("Endpoint not found"))
}
