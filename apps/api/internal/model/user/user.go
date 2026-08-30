package user

import "github.com/chandankrr/loreline/internal/model"

type User struct {
	model.Base
	Email         string  `json:"email" db:"email"`
	Name          string  `json:"name" db:"name"`
	EmailVerified bool    `json:"emailVerified" db:"email_verified"`
	Image         *string `json:"image" db:"image"`
}
