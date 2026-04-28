package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/common/types"
	"github.com/ayopedro/seazus-go/internal/models"
)

func (h *handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.LoginUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, appErrors.ErrInvalidInput)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		utils.WriteError(w, r, appErrors.ErrInvalidInput)
		return
	}

	authUser, err := h.service.Auth.LoginUser(r.Context(), payload)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidCredentials) {
			utils.WriteError(w, r, err)
			return
		}
		utils.WriteError(w, r, err)
		return
	}

	setCookie(w, "auth_token", authUser.Token, time.Now().Add(1*time.Hour))

	response := types.APIResponseBody{
		Status:  true,
		Message: "Login successful",
		Data:    authUser,
	}

	utils.WriteJSON(w, r, http.StatusOK, response)
}

func setCookie(w http.ResponseWriter, name, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	payload := &models.CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		utils.WriteError(w, r, appErrors.ErrInvalidInput)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		utils.WriteError(w, r, appErrors.ErrInvalidInput)
		return
	}

	err := h.service.Auth.CreateUser(r.Context(), payload)

	if err != nil {
		if errors.Is(err, appErrors.ErrConflict) {
			utils.WriteError(w, r, err)
			return
		}
		utils.WriteError(w, r, err)
		return
	}
	response := types.APIResponseBody{
		Status:  true,
		Message: "User created successfully",
	}

	utils.WriteJSON(w, r, http.StatusCreated, response)
}
