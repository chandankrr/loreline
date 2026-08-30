package v1

import (
	"github.com/chandankrr/loreline/internal/handler"
	"github.com/chandankrr/loreline/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {
	// Register auth routes
	registerAuthRoutes(router, handlers.Auth, middleware.Auth)
}
