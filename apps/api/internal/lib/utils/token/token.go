package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const TokenBytes = 32

// Generate generates a cryptographically secure random token
func Generate() (string, error) {
	b := make([]byte, TokenBytes)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Hash hashes a token before database storage/lookup
func Hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
