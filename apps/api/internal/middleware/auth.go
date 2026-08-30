package middleware

import (
	"strings"
	"time"

	"github.com/chandankrr/loreline/internal/errs"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/chandankrr/loreline/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	server      *server.Server
	authService *service.AuthService
}

func NewAuthMiddleware(s *server.Server, authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		server:      s,
		authService: authService,
	}
}

func (auth *AuthMiddleware) RequiredAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			auth.server.Logger.Error().
				Str("function", "RequiredAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("missing authorization header")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		// Expected format: Authorization: Bearer <access_token>
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			auth.server.Logger.Error().
				Str("function", "RequiredAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("invalid authorization header format")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		accessToken := parts[1]

		if accessToken == "" {
			auth.server.Logger.Error().
				Str("function", "RequiredAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("empty access token")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		claims, err := auth.authService.ValidateToken(accessToken)
		if err != nil {
			auth.server.Logger.Error().
				Err(err).
				Str("function", "RequiredAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("invalid or expired access token")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		// Extract user ID from JWT claims
		userID := claims.Subject
		if userID == "" {
			auth.server.Logger.Error().
				Str("function", "RequiredAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("user id missing from token claims")

			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		c.Set("user_id", userID)

		auth.server.Logger.Info().
			Str("function", "RequiredAuth").
			Str("user_id", userID).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)
	}
}
