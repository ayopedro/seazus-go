package main

import (
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/handlers"
	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/middleware"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
)

type application struct {
	config  *config.Config
	limiter ratelimiter.Limiter
	h       handlers.Handler
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// =======================
	// ROUTES				||
	// =======================
	mux.HandleFunc("GET /v1", app.h.IndexHandler)
	// Health route
	mux.HandleFunc("GET /v1/health", app.h.HealthCheckHandler)

	// Auth routes
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /login", app.h.LoginHandler)
	authMux.HandleFunc("POST /register", app.h.RegisterHandler)
	// authMux.HandleFunc("POST /forgot-password", nil)

	// protected route
	usersMux := http.NewServeMux()
	urlsMux := http.NewServeMux()

	// user related route
	usersMux.HandleFunc("GET /me", app.h.GetMyProfileHandler)
	usersMux.HandleFunc("GET /urls", app.h.GetUserURLSHandler)

	urlsMux.HandleFunc("GET /{id}", app.h.GetURLByIdHandler)
	urlsMux.HandleFunc("POST /", app.h.CreateURLHandler)

	// Route grouping
	mux.Handle("/v1/auth/", http.StripPrefix("/v1/auth", authMux))
	mux.Handle("/v1/users/", http.StripPrefix("/v1/users", app.h.Protected(usersMux)))
	mux.Handle("/v1/urls/", http.StripPrefix("/v1/urls", app.h.Protected(urlsMux)))

	mux.HandleFunc("GET /{short_url}", app.h.ShortURLRedirectHandler)

	// Not found handler
	mux.HandleFunc("/", app.h.NotFoundHandler)

	var h http.Handler = mux
	h = middleware.RequestLogger(h)
	h = middleware.CORS(app.config.TrustedOrigins)(h)
	if app.config.AppEnv == "production" {
		h = middleware.RateLimiter(app.limiter)(h)
	}
	h = middleware.RecoverPanic(h)
	return h
}
