package account

import (
	"time"

	"github.com/chandankrr/loreline/internal/model"
	"github.com/google/uuid"
)

type Provider string

const (
	ProviderCredential Provider = "credential"
	ProviderGoogle     Provider = "google"
)

type Account struct {
	model.Base
	AccountID             string     `json:"accountId" db:"account_id"`
	ProviderID            Provider   `json:"providerId" db:"provider_id"`
	UserID                uuid.UUID  `json:"userId" db:"user_id"`
	AccessToken           *string    `json:"accessToken" db:"access_token"`
	RefreshToken          *string    `json:"refreshToken" db:"refresh_token"`
	IDToken               *string    `json:"idToken" db:"id_token"`
	AccessTokenExpiresAt  *time.Time `json:"accessTokenExpiresAt" db:"access_token_expires_at"`
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt" db:"refresh_token_expires_at"`
	Scope                 *string    `json:"scope" db:"scope"`
	Password              *string    `json:"password" db:"password"`
}
