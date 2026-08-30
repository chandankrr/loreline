package handler

import (
	"github.com/chandankrr/loreline/internal/server"
	"github.com/chandankrr/loreline/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Auth    *AuthHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Auth:    NewAuthHandler(s, services.Auth),
	}
}
