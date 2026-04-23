package main

import (
	"net/http"

	"github.com/ayopedro/seazus-go/lib/config"
)

type application struct {
	config   *config.Config
	dbConfig dbConfig
}

type dbConfig struct {
	DBURL        string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// =======================
	// ROUTES				||
	// =======================
	// Health route
	mux.HandleFunc("GET /v1/health", app.HealthCheckHandler)

	handler := logRequest(recoverPanic(mux))
	return handler
}
