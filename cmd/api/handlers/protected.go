package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ayopedro/seazus-go/cmd/api/response"
	"github.com/ayopedro/seazus-go/internal/apperrors"
)

type contextKey string

const userContextKey = contextKey("user")

func (h *Handler) Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		if cookie, err := r.Cookie("auth_token"); err == nil {
			token = cookie.Value
		}

		if token == "" {
			authHeader := r.Header.Get("Authorization")

			if authToken, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
				token = authToken
			}
		}

		if token == "" {
			response.WriteError(w, http.StatusUnauthorized, apperrors.ErrUnauthorized)
			return
		}

		claims, err := h.authValidator.Validate(token)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, apperrors.ErrUnauthorized)
			return
		}

		user, err := h.user.GetUserProfile(r.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, apperrors.ErrUserNotFound) {
				response.WriteError(w, http.StatusNotFound, apperrors.ErrUserNotFound)
			} else {
				response.WriteError(w, http.StatusForbidden, apperrors.ErrForbidden)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
