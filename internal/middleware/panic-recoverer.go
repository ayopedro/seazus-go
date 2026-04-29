package middleware

import (
	"errors"
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/response"
	"github.com/ayopedro/seazus-go/internal/logger"
	"go.uber.org/zap"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", r.URL.Path),
				)

				response.WriteError(
					w,
					errors.New(http.StatusText(http.StatusInternalServerError)),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
