package verification

import (
	"time"

	"github.com/chandankrr/loreline/internal/model"
)

type Verification struct {
	model.Base
	Identifier string    `json:"identifier" db:"identifier"`
	Value      string    `json:"value" db:"value"`
	ExpiresAt  time.Time `json:"expiresAt" db:"expires_at"`
}
