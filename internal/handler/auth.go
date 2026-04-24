package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ayopedro/seazus-go/internal/models"
	"github.com/ayopedro/seazus-go/internal/types"
	"github.com/ayopedro/seazus-go/internal/utils"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.LoginUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, models.ErrInvalidPayload)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		utils.WriteJSON(w, r, http.StatusBadRequest, types.APIResponseBody{
			Status:  false,
			Message: "email and password are required",
		})
		return
	}

	authUser, err := h.AuthService.LoginUser(r.Context(), payload)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
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

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, http.StatusBadRequest, models.ErrInvalidPayload)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		utils.WriteJSON(w, r, http.StatusBadRequest, types.APIResponseBody{
			Status:  false,
			Message: "email and password are required",
		})
		return
	}

	err := h.AuthService.CreateUser(r.Context(), payload)

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
