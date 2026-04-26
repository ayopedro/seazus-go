package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	utils "github.com/ayopedro/seazus-go/internal/common"
	appErrors "github.com/ayopedro/seazus-go/internal/common/app_errors"
	"github.com/ayopedro/seazus-go/internal/logger"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
	"go.uber.org/zap"
)

type contextKey string

const userContextKey = contextKey("user")

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.Info("http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", r.URL.Path),
				)

				utils.WriteError(
					w,
					r,
					http.StatusInternalServerError,
					errors.New(http.StatusText(http.StatusInternalServerError)),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RateLimiter(limiter ratelimiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			if host, _, err := net.SplitHostPort(ip); err != nil {
				ip = host
			}

			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ips := strings.Split(xff, ",")
				ip = strings.TrimSpace(ips[0])
			}

			if allow, retryAfter := limiter.Allow(ip); !allow {
				w.Header().Set("Retry-After", fmt.Sprintf("%.f", retryAfter.Seconds()))
				utils.WriteError(w, r, http.StatusTooManyRequests, errors.New(http.StatusText(429)))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CORS(trustedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")

			if origin != "" {
				for _, trustedOrigin := range trustedOrigins {
					if origin == trustedOrigin {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

							w.WriteHeader(http.StatusOK)
							return
						}
						break
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

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
