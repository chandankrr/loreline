package errs

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrAccountNotFound      = errors.New("account not found")
	ErrSessionNotFound      = errors.New("session not found")
	ErrVerificationNotFound = errors.New("verification record not found")
)
