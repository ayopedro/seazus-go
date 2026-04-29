package handlers

import (
	"github.com/ayopedro/seazus-go/internal/service/auth"
	"github.com/ayopedro/seazus-go/internal/service/url"
	"github.com/ayopedro/seazus-go/internal/service/user"
)

type Handler struct {
	auth          auth.Service
	authValidator auth.TokenValidator
	user          user.Service
	url           url.Service
}

func NewHandler(auth auth.Service, validator auth.TokenValidator, user user.Service, url url.Service) *Handler {
	return &Handler{auth, validator, user, url}
}
