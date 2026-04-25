package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/repository"
	"github.com/ayopedro/seazus-go/internal/service"
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
	GetMyProfile(w http.ResponseWriter, r *http.Request)
	GetURLById(w http.ResponseWriter, r *http.Request)
	GetUserURLS(w http.ResponseWriter, r *http.Request)
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)

	Protected(next http.Handler) http.Handler
}

func NewHandler(
	c *config.Config,
	db *sql.DB,
) Handler {
	// repositories
	ur := repository.NewUserRepository(db)
	urlr := repository.NewURLRepository(db)

	// services
	us := service.NewUserService(ur, urlr)
	as := service.NewAuthService(ur, c.JWTSecret)
	urls := service.NewURLService(urlr)

	return &handler{
		config:      c,
		authService: as,
		userService: us,
		urlService:  urls,
	}
}
