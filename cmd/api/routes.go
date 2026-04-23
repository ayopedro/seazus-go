package main

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/handler"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
)

type application struct {
	config  *config.Config
	db      *sql.DB
	limiter ratelimiter.Limiter
	wg      sync.WaitGroup
	handler *handler.Handler
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// =======================
	// ROUTES				||
	// =======================
	// Health route
	mux.HandleFunc("GET /v1/health", app.handler.HealthCheckHandler)

	var h http.Handler = mux
	h = handler.RecoverPanic(h)
	h = handler.LogRequest(h)
	h = handler.CORS(app.config.TrustedOrigins)(h)
	if app.config.AppEnv == "production" {
		h = handler.RateLimiter(app.limiter)(h)
	}
	return h
}
