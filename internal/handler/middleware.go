package handler

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ayopedro/seazus-go/internal/logger"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
	"github.com/ayopedro/seazus-go/internal/utils"
	"go.uber.org/zap"
)

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

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
