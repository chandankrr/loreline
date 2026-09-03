package oauth

import (
	"net/http"

	"github.com/chandankrr/loreline/internal/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const sessionMaxAgeSeconds = 10 * 60 // 10 minutes; only needs to survive the OAuth round trip

func Setup(cfg *config.Config) {
	store := sessions.NewCookieStore([]byte(cfg.OAuth.SessionSecret))
	store.MaxAge(sessionMaxAgeSeconds)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = cfg.Primary.Env == "production"
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

	goth.UseProviders(
		google.New(
			cfg.OAuth.Google.ClientID,
			cfg.OAuth.Google.ClientSecret,
			callbackURL(cfg, "google"),
			"email", "profile",
		),
	)
}

func callbackURL(cfg *config.Config, provider string) string {
	return cfg.OAuth.CallbackURL + "/" + provider + "/callback"
}
