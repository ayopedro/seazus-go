package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/models"
)

func (h *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.LoginUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, appErrors.ErrInvalidInput)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		utils.WriteError(w, r, http.StatusBadRequest, appErrors.ErrInvalidInput)
		return
	}

	authUser, err := h.authService.LoginUser(r.Context(), payload)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidCredentials) {
			utils.WriteError(w, r, http.StatusUnauthorized, err)
			return
		}
		utils.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response := types.APIResponseBody{
		Status:  true,
		Message: "Login successful",
		Data:    authUser,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func (h *handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, appErrors.ErrInvalidInput)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		utils.WriteJSON(w, r, http.StatusBadRequest, types.APIResponseBody{
			Status:  false,
			Message: "email and password are required",
		})
		return
	}

	err := h.authService.CreateUser(r.Context(), payload)

	if err != nil {
		if errors.Is(err, appErrors.ErrConflict) {
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
