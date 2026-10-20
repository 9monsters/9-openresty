package auth

import (
	"encoding/base64"
	"crypto/rand"
	"fmt"
)

// GenerateKey produces a random 32-byte signing key, encoded as base64.
func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("generate key: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(key), nil
}