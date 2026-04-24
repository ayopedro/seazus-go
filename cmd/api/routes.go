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
	h       *handler.Handler
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// =======================
	// ROUTES				||
	// =======================
	// Health route
	mux.HandleFunc("GET /v1/health", app.h.HealthCheckHandler)
	mux.HandleFunc("GET /v1", app.h.IndexHandler)

	// Auth routes
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /login", app.h.LoginHandler)
	authMux.HandleFunc("POST /register", app.h.RegisterHandler)
	// authMux.HandleFunc("POST /forgot-password", nil)

	// Route grouping
	mux.Handle("/v1/auth/", http.StripPrefix("/v1/auth", authMux))

	var h http.Handler = mux
	h = handler.LogRequest(h)
	h = handler.CORS(app.config.TrustedOrigins)(h)
	if app.config.AppEnv == "production" {
		h = handler.RateLimiter(app.limiter)(h)
	}
	h = handler.RecoverPanic(h)
	return h
}
