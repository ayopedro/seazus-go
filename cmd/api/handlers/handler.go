package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/repository"
	"github.com/ayopedro/seazus-go/internal/service"
	"go.uber.org/zap"
)

type handler struct {
	config  *config.Config
	service *service.Service
}

type Handler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	GetMyProfileHandler(w http.ResponseWriter, r *http.Request)
	GetURLByIdHandler(w http.ResponseWriter, r *http.Request)
	GetUserURLSHandler(w http.ResponseWriter, r *http.Request)
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)
	CreateURLHandler(w http.ResponseWriter, r *http.Request)

	Protected(next http.Handler) http.Handler
}

func NewHandler(cfg *config.Config, db *sql.DB, logger *zap.Logger) Handler {
	// repositories
	repo := repository.NewRepository(db, logger)

	// services
	s := service.NewService(repo, logger, cfg.JWTSecret)

	return &handler{
		config:  cfg,
		service: s,
	}
}
