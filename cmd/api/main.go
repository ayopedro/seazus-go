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

	"github.com/ayopedro/seazus-go/lib/config"
	"github.com/ayopedro/seazus-go/lib/logger"
	"go.uber.org/zap"
)

var cfg *config.Config

func createServer(cfg *config.Config, handler http.Handler) *http.Server {
	addr := fmt.Sprintf(":%s", cfg.Port)

	logger.Info("Configuring server settings", zap.String("addr", addr))

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
	cfg = config.Load()
	logger.Init(cfg.AppEnv)
	defer logger.Sync()

	app := &application{
		config: cfg,
		dbConfig: dbConfig{
			DBURL:        cfg.DBURL,
			maxOpenConns: cfg.DBMaxOpenConns,
			maxIdleConns: cfg.DBMaxIdleConns,
			maxIdleTime:  cfg.DBMaxIdleTime,
		},
	}
	server := createServer(cfg, app.routes())
	runServer(server)
}
