package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/db"
	"github.com/ayopedro/seazus-go/internal/handler"
	"github.com/ayopedro/seazus-go/internal/logger"
	ratelimiter "github.com/ayopedro/seazus-go/internal/middleware"
	"go.uber.org/zap"
)

func (app *application) createServer(handler http.Handler) *http.Server {
	addr := fmt.Sprintf(":%s", app.config.Port)

	logger.Info("Configuring server settings")

	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}
}

func runServer(srv *http.Server) {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Info("Shutting down server", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.Info("Starting server",
		zap.String("addr", srv.Addr),
	)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Server failed to start", zap.Error(err))
	}

	err = <-shutdownError
	if err != nil {
		logger.Error("Graceful shutdown error", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)
	defer logger.Sync()

	limiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.Limiter.RequestsPerTimeframe,
		cfg.Limiter.Timeframe,
	)

	db, err := db.New(cfg.DB)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Database is connected")

	h := &handler.Handler{
		AppConfig: cfg,
		DB:        db,
	}

	app := &application{
		config:  cfg,
		db:      db,
		limiter: limiter,
		handler: h,
	}

	server := app.createServer(app.routes())
	runServer(server)
}
