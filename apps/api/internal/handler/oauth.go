package handler

import (
	"net/http"

	"github.com/chandankrr/loreline/internal/logger"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/chandankrr/loreline/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type OAuthHandler struct {
	Handler
	authService *service.AuthService
}

func NewOAuthHandler(s *server.Server, authService *service.AuthService) *OAuthHandler {
	return &OAuthHandler{
		Handler:     NewHandler(s),
		authService: authService,
	}
}

func (h *OAuthHandler) BeginAuth(c echo.Context) error {
	req := withProviderParam(c)
	gothic.BeginAuthHandler(c.Response(), req)
	return nil
}

func (h *OAuthHandler) Callback(c echo.Context) error {
	logger := logger.GetLogger(c)
	provider := c.Param("provider")

	req := withProviderParam(c)

	frontendURL := h.server.Config.OAuth.FrontendURL

	gothUser, err := gothic.CompleteUserAuth(c.Response(), req)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("provider", provider).
			Msg("oauth callback failed")
		return c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/sign-in?error=oauth_failed")
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	_, refreshToken, err := h.authService.OAuthLogin(c, gothUser, ipAddress, userAgent)
	if err != nil {
		logger.Error().
			Err(err).
			Str("provider", provider).
			Msg("oauth login failed")
		return c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/sign-in?error=oauth_failed")
	}

	h.setRefreshTokenCookie(c, refreshToken)

	return c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/callback")
}

// withProviderParam bridges echo's :provider path param to gothic
func withProviderParam(c echo.Context) *http.Request {
	req := c.Request()

	q := req.URL.Query()
	q.Set("provider", c.Param("provider"))
	req.URL.RawQuery = q.Encode()

	return req
}
