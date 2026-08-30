package v1

import (
	"github.com/chandankrr/loreline/internal/handler"
	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(r *echo.Group, h *handler.AuthHandler) {
	auth := r.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.RefreshToken)
}
