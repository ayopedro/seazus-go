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
	config      *config.Config
	authService service.AuthService
	userService service.UserService
	urlService  service.URLService
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

func NewHandler(
	c *config.Config,
	db *sql.DB,
	l *zap.Logger,
) Handler {
	// repositories
	ur := repository.NewUserRepository(db, l)
	urlr := repository.NewURLRepository(db, l)

	// services
	us := service.NewUserService(ur, urlr, l)
	as := service.NewAuthService(ur, c.JWTSecret, l)
	urls := service.NewURLService(urlr, l)

	return &handler{
		config:      c,
		authService: as,
		userService: us,
		urlService:  urls,
	}
}
