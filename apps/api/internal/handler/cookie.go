package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const refreshTokenCookieName = "refresh_token"

const refreshTokenCookiePath = "/api/v1/auth"

func (h Handler) setRefreshTokenCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     refreshTokenCookiePath,
		Expires:  time.Now().Add(h.server.Config.Auth.RefreshTokenTTL),
		HttpOnly: true,
		Secure:   h.server.Config.Primary.Env == "production",
		SameSite: http.SameSiteStrictMode,
	})
}

func (h Handler) clearRefreshTokenCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     refreshTokenCookiePath,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.server.Config.Primary.Env == "production",
		SameSite: http.SameSiteStrictMode,
	})
}
