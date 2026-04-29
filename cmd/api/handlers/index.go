package handlers

import (
	"net/http"
	"time"

	"github.com/ayopedro/seazus-go/cmd/api/response"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	result := response.APIResponseBody{
		Status:  true,
		Message: "Welcome to the Seazus(Go) API",
		Data:    []string{},
	}

	response.WriteJSON(w, r, http.StatusOK, result)
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	result := response.APIResponseBody{
		Status:  true,
		Message: "Server is healthy",
		Data: map[string]string{
			"status": "available",
			"time":   time.Now().String(),
		},
	}

	response.WriteJSON(w, r, http.StatusOK, result)
}

func (h *Handler) ShortURLRedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("short_url")
	original_url, err := h.url.GetOriginalURL(r.Context(), path)

	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	http.Redirect(w, r, original_url, http.StatusSeeOther)
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	result := response.APIResponseBody{
		Status:  false,
		Message: "Endpoint not found",
		Data:    []string{},
	}

	response.WriteJSON(w, r, http.StatusNotFound, result)
}
