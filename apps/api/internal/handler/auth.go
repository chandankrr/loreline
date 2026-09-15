package handler

import (
	"errors"
	"net/http"

	"github.com/chandankrr/loreline/internal/dto"
	"github.com/chandankrr/loreline/internal/errs"
	"github.com/chandankrr/loreline/internal/middleware"
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
		func(c echo.Context, payload *dto.RegisterPayload) (*dto.MessageResponse, error) {
			_, err := h.authService.Register(c, payload)
			if err != nil {
				if errors.Is(err, service.ErrEmailInUse) {
					code := "EMAIL_ALREADY_IN_USE"
					return nil, errs.NewConflictError("Email already in use", false, &code)
				}
				return nil, err
			}

			return &dto.MessageResponse{Message: "User registered successfully"}, nil
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
				switch {
				case errors.Is(err, service.ErrEmailNotVerified):
					code := "EMAIL_NOT_VERIFIED"
					return nil,
						errs.NewUnauthorizedError("Email not verified", false, &code)
				case errors.Is(err, service.ErrInvalidCredentials):
					return nil,
						errs.NewUnauthorizedError("Invalid email or password", false, nil)
				default:
					return nil, err
				}
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
					errs.NewUnauthorizedError("Missing refresh token", false, nil)
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
					h.clearRefreshTokenCookie(c)
					return nil,
						errs.NewUnauthorizedError("Invalid token", false, nil)

				case errors.Is(err, service.ErrExpiredToken):
					h.clearRefreshTokenCookie(c)
					return nil,
						errs.NewUnauthorizedError("Token has expired", false, nil)

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

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.VerifyEmailPayload) (*dto.MessageResponse, error) {
			err := h.authService.VerifyEmail(c, payload.Email, payload.Code)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrInvalidVerificationCode):
					code := "INVALID_CODE"
					return nil,
						errs.NewBadRequestError("Invalid verification code", false, &code, nil, nil)
				case errors.Is(err, service.ErrVerificationExpired):
					code := "CODE_EXPIRED"
					return nil,
						errs.NewBadRequestError("Verification code has expired", false, &code, nil, nil)
				case errors.Is(err, service.ErrEmailAlreadyVerified):
					return nil,
						errs.NewBadRequestError("Email already verified", false, nil, nil, nil)
				case errors.Is(err, service.ErrVerificationAttemptsExceeded):
					return nil,
						errs.NewTooManyRequestsError("Too many verification attempts. Please request a new code", false)
				default:
					return nil, err
				}
			}

			return &dto.MessageResponse{Message: "Email verified successfully"}, nil
		},
		http.StatusOK,
		&dto.VerifyEmailPayload{},
	)(c)
}

func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.ResendVerificationPayload) (*dto.MessageResponse, error) {
			err := h.authService.ResendVerificationEmail(c, payload.Email)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrEmailAlreadyVerified):
					return nil,
						errs.NewBadRequestError("Email already verified", false, nil, nil, nil)
				case errors.Is(err, service.ErrTooManyRequests):
					return nil,
						errs.NewTooManyRequestsError("Please wait before requesting another code", false)
				default:
					return nil, err
				}
			}

			return &dto.MessageResponse{
				Message: "If your email is registered and unverified, a verification code has been sent",
			}, nil
		},
		http.StatusOK,
		&dto.ResendVerificationPayload{},
	)(c)
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.ForgotPasswordPayload) (*dto.MessageResponse, error) {
			err := h.authService.ForgotPassword(c, payload.Email)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrTooManyRequests):
					return nil,
						errs.NewTooManyRequestsError("Please wait before requesting another password reset", false)
				default:
					return nil, err
				}
			}

			return &dto.MessageResponse{
				Message: "If your email is registered, a password reset link has been sent",
			}, nil
		},
		http.StatusOK,
		&dto.ForgotPasswordPayload{},
	)(c)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.ResetPasswordPayload) (*dto.MessageResponse, error) {
			err := h.authService.ResetPassword(c, payload.Token, payload.Password)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrInvalidResetToken):
					return nil,
						errs.NewBadRequestError("Invalid or expired reset token", false, nil, nil, nil)
				default:
					return nil, err
				}
			}

			return &dto.MessageResponse{
				Message: "Password reset successfully",
			}, nil
		},
		http.StatusOK,
		&dto.ResetPasswordPayload{},
	)(c)
}

func (h *AuthHandler) Me(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.EmptyPayload) (*dto.UserResponse, error) {
			userID := middleware.GetUserID(c)

			user, err := h.authService.GetCurrentUser(c, userID)
			if err != nil {
				return nil, err
			}

			return dto.ToUserResponse(user), nil
		},
		http.StatusOK,
		&dto.EmptyPayload{},
	)(c)
}
