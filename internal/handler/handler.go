package handler

import (
	"github.com/ayopedro/seazus-go/internal/config"
	"github.com/ayopedro/seazus-go/internal/service"
)

type Handler struct {
	AppConfig   *config.Config
	UserService service.UserService
}
