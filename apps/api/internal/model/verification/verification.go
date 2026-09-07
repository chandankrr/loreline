package verification

import (
	"time"

	"github.com/chandankrr/loreline/internal/model"
)

type VerificationType string

const (
	TypeEmailVerification VerificationType = "email_verification"
	TypePasswordReset     VerificationType = "password_reset"
)

type Verification struct {
	model.Base
	Identifier string           `json:"identifier" db:"identifier"`
	Type       VerificationType `json:"type" db:"type"`
	Value      string           `json:"value" db:"value"`
	ExpiresAt  time.Time        `json:"expiresAt" db:"expires_at"`
	Attempts   int              `json:"attempts" db:"attempts"`
}
