package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ayopedro/seazus-go/cmd/api/dto"
	"github.com/ayopedro/seazus-go/cmd/api/response"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/models"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	payload := &dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.WriteError(w, appErrors.ErrInvalidInput)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		response.WriteError(w, appErrors.ErrInvalidInput)
		return
	}

	authUser, err := h.auth.LoginUser(
		r.Context(),
		&models.LoginUser{
			Email:    payload.Email,
			Password: payload.Password,
		},
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidCredentials) {
			response.WriteError(w, err)
			return
		}
		response.WriteError(w, err)
		return
	}

	setCookie(w, "auth_token", authUser.Token, time.Now().Add(1*time.Hour))

	data := dto.LoginResponse{
		Token: authUser.Token,
		User: dto.AuthUser{
			Id:         authUser.User.Id,
			FirstName:  authUser.User.FirstName,
			Email:      authUser.User.Email,
			IsVerified: authUser.User.IsVerified,
		},
	}

	response.WriteJSON(w, http.StatusOK, "Login successful", data)
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

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	payload := &dto.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.WriteError(w, appErrors.ErrInvalidInput)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		response.WriteError(w, appErrors.ErrInvalidInput)
		return
	}

	err := h.auth.CreateUser(
		r.Context(),
		&models.CreateUser{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Password:  payload.Password,
		},
	)

	if err != nil {
		if errors.Is(err, appErrors.ErrConflict) {
			response.WriteError(w, err)
			return
		}
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, "User created successfully", nil)
}
