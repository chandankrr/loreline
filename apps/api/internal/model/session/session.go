package session

import (
	"net"
	"time"

	"github.com/chandankrr/loreline/internal/model"
	"github.com/google/uuid"
)

type Session struct {
	model.Base
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
	IPAddress *net.IP   `json:"ipAddress" db:"ip_address"`
	UserAgent *string   `json:"userAgent" db:"user_agent"`
}
