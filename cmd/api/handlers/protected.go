package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
)

type contextKey string

const userContextKey = contextKey("user")

func (h *handler) Protected(next http.Handler) http.Handler {
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
			utils.WriteError(w, r, http.StatusUnauthorized, appErrors.ErrInvalidToken)
			return
		}

		claims, err := utils.ValidateToken(token, h.config.JWTSecret)
		if err != nil {
			utils.WriteError(w, r, http.StatusUnauthorized, appErrors.ErrInvalidToken)
			return
		}

		user, err := h.userService.GetUserProfile(r.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				utils.WriteError(w, r, http.StatusUnauthorized, appErrors.ErrNotFound)
			} else {
				utils.WriteError(w, r, http.StatusInternalServerError, appErrors.ErrForbidden)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
