package main

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/ayopedro/seazus-go/cmd/api/handlers"
	"github.com/ayopedro/seazus-go/internal/config"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
)

type application struct {
	config  *config.Config
	db      *sql.DB
	limiter ratelimiter.Limiter
	wg      sync.WaitGroup
	h       handlers.Handler
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

	// protected route
	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /me", app.h.GetMyProfile)
	userMux.HandleFunc("GET /urls/{id}", app.h.GetURLById)

	// Route grouping
	mux.Handle("/v1/auth/", http.StripPrefix("/v1/auth", authMux))
	mux.Handle("/v1/users/", http.StripPrefix("/v1/users", app.h.Protected(userMux)))

	var h http.Handler = mux
	h = handlers.LogRequest(h)
	h = handlers.CORS(app.config.TrustedOrigins)(h)
	if app.config.AppEnv == "production" {
		h = handlers.RateLimiter(app.limiter)(h)
	}
	h = handlers.RecoverPanic(h)
	return h
}
