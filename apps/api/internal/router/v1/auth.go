package v1

import (
	"github.com/chandankrr/loreline/internal/handler"
	"github.com/chandankrr/loreline/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(
	r *echo.Group,
	h *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	m *middleware.AuthMiddleware,
) {
	auth := r.Group("/auth")

	// Credential
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout, m.RequiredAuth)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/email/verify", h.VerifyEmail)
	auth.POST("/email/resend", h.ResendVerificationEmail)
	auth.POST("/password/forgot", h.ForgotPassword)
	auth.POST("/password/reset", h.ResetPassword)

	// OAuth
	auth.GET("/:provider", oauthHandler.BeginAuth)
	auth.GET("/:provider/callback", oauthHandler.Callback)
}
