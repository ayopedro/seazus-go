package handlers

import (
	"net/http"
	"time"

	utils "github.com/ayopedro/seazus-go/internal/common"
	"github.com/ayopedro/seazus-go/internal/common/types"
)

func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  true,
		Message: "Welcome to the Seazus(Go) API",
		Data:    []string{},
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func (h *handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  true,
		Message: "Server is healthy",
		Data: map[string]string{
			"status":      "available",
			"environment": h.config.AppEnv,
			"time":        time.Now().String(),
		},
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func (h *handler) ShortURLRedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("short_url")
	original_url, err := h.service.URL.GetOriginalURL(r.Context(), path)

	if err != nil {
		utils.WriteError(w, r, err)
		return
	}

	http.Redirect(w, r, original_url, http.StatusSeeOther)
}

func (h *handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := types.APIResponseBody{
		Status:  false,
		Message: "Endpoint not found",
		Data:    []string{},
	}

	utils.WriteJSON(w, r, http.StatusNotFound, response)
}
