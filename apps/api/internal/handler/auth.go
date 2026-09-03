package handler

import (
	"errors"
	"net/http"

	"github.com/chandankrr/loreline/internal/dto"
	"github.com/chandankrr/loreline/internal/errs"
	"github.com/chandankrr/loreline/internal/model/user"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/chandankrr/loreline/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Handler
	authService *service.AuthService
}

func NewAuthHandler(s *server.Server, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Handler:     NewHandler(s),
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.RegisterPayload) (*user.User, error) {
			user, err := h.authService.Register(c, payload)
			if err != nil {
				if errors.Is(err, service.ErrEmailInUse) {
					code := "EMAIL_ALREADY_IN_USE"
					return nil, errs.NewConflictError("Email already in use", false, &code)
				}
				return nil, err
			}

			return user, nil
		},
		http.StatusCreated,
		&dto.RegisterPayload{},
	)(c)
}

func (h *AuthHandler) Login(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.LoginPayload) (*dto.LoginResponse, error) {
			ipAddress := c.RealIP()
			userAgent := c.Request().UserAgent()

			accessToken, refreshToken, err := h.authService.Login(
				c,
				payload,
				ipAddress,
				userAgent,
			)
			if err != nil {
				if errors.Is(err, service.ErrInvalidCredentials) {
					return nil,
						errs.NewUnauthorizedError("Invalid email or password", false)
				}
				return nil, err
			}

			h.setRefreshTokenCookie(c, refreshToken)

			return &dto.LoginResponse{AccessToken: accessToken}, nil
		},
		http.StatusOK,
		&dto.LoginPayload{},
	)(c)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *dto.EmptyPayload) error {
			cookie, err := c.Cookie(refreshTokenCookieName)
			if err == nil && cookie.Value != "" {
				if err := h.authService.Logout(c, cookie.Value); err != nil {
					return err
				}
			}

			h.clearRefreshTokenCookie(c)
			return nil
		},
		http.StatusNoContent,
		&dto.EmptyPayload{},
	)(c)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.EmptyPayload) (*dto.RefreshResponse, error) {
			cookie, err := c.Cookie(refreshTokenCookieName)
			if err != nil || cookie.Value == "" {
				return nil,
					errs.NewUnauthorizedError("Missing refresh token", false)
			}

			ipAddress := c.RealIP()
			userAgent := c.Request().UserAgent()

			accessToken, refreshToken, err := h.authService.RefreshAccessToken(
				c,
				cookie.Value,
				ipAddress,
				userAgent,
			)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrInvalidToken):
					return nil,
						errs.NewUnauthorizedError("Invalid token", false)

				case errors.Is(err, service.ErrExpiredToken):
					return nil,
						errs.NewUnauthorizedError("Token has expired", false)

				default:
					return nil, err
				}
			}

			h.setRefreshTokenCookie(c, refreshToken)

			return &dto.RefreshResponse{AccessToken: accessToken}, nil
		},
		http.StatusOK,
		&dto.EmptyPayload{},
	)(c)
}
