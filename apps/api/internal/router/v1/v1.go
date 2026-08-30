package v1

import (
	"github.com/chandankrr/loreline/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers) {
	registerAuthRoutes(router, handlers.Auth)
}
