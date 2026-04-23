package main

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/ayopedro/seazus-go/internals/ratelimiter"
	"github.com/ayopedro/seazus-go/lib/config"
)

type application struct {
	config  *config.Config
	db      *sql.DB
	limiter ratelimiter.Limiter
	wg      sync.WaitGroup
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// =======================
	// ROUTES				||
	// =======================
	// Health route
	mux.HandleFunc("GET /v1/health", app.HealthCheckHandler)

	var handler http.Handler = mux
	handler = logRequest(recoverPanic(mux))
	handler = app.cors(handler)

	if app.config.AppEnv == "production" {
		handler = app.rateLimiter(handler)
	}
	return handler
}
