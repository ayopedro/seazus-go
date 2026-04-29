package main

import (
	"net/http"

	"github.com/ayopedro/seazus-go/cmd/api/handlers"
	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/middleware"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
)

func routes(cfg *config.Config, h *handlers.Handler, limiter ratelimiter.Limiter) http.Handler {
	mux := http.NewServeMux()

	// =======================
	// ROUTES				||
	// =======================
	mux.HandleFunc("GET /v1", h.IndexHandler)
	// Health route
	mux.HandleFunc("GET /v1/health", h.HealthCheckHandler)

	// Auth routes
	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /login", h.LoginHandler)
	authMux.HandleFunc("POST /register", h.RegisterHandler)
	// authMux.HandleFunc("POST /forgot-password", nil)

	// protected route
	usersMux := http.NewServeMux()
	urlsMux := http.NewServeMux()

	// user related route
	usersMux.HandleFunc("GET /me", h.GetMyProfileHandler)
	usersMux.HandleFunc("GET /urls", h.GetUserURLSHandler)

	urlsMux.HandleFunc("GET /{id}", h.GetURLByIdHandler)
	urlsMux.HandleFunc("POST /", h.CreateURLHandler)
	urlsMux.HandleFunc("PATCH /{id}", h.UpdateURLHandler)
	urlsMux.HandleFunc("DELETE /{id}", h.DeleteURLHandler)

	// Route grouping
	mux.Handle("/v1/auth/", http.StripPrefix("/v1/auth", authMux))
	mux.Handle("/v1/users/", http.StripPrefix("/v1/users", h.Protected(usersMux)))
	mux.Handle("/v1/urls/", http.StripPrefix("/v1/urls", h.Protected(urlsMux)))

	mux.HandleFunc("GET /{short_url}", h.ShortURLRedirectHandler)

	// Not found handler
	mux.HandleFunc("/", h.NotFoundHandler)

	var handler http.Handler = mux
	handler = middleware.RequestLogger(handler)
	handler = middleware.CORS(cfg.TrustedOrigins)(handler)
	if cfg.AppEnv == "production" {
		handler = middleware.RateLimiter(limiter)(handler)
	}
	handler = middleware.RecoverPanic(handler)
	return handler
}
