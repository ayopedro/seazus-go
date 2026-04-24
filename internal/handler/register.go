package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/types"
	"github.com/ayopedro/seazus-go/internal/utils"
)

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	up := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(up); err != nil {
		response := types.APIResponseBody{
			Status:  false,
			Message: "Invalid request body",
		}
		utils.WriteJSON(w, r, http.StatusBadRequest, response)
		return
	}

	if up.Email == "" || up.Password == "" {
		utils.WriteJSON(w, r, http.StatusBadRequest, types.APIResponseBody{
			Status:  false,
			Message: "email and password are required",
		})
		return
	}

	err := h.UserService.CreateUser(r.Context(), up)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			utils.WriteError(w, r, http.StatusConflict, err)
			return
		}
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	response := types.APIResponseBody{
		Status:  true,
		Message: "User created successfully",
	}

	utils.WriteJSON(w, r, http.StatusCreated, response)
}
