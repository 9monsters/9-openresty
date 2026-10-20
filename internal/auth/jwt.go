// Package auth provides JWT authentication for the web UI.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Claims is the JWT payload.
type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"un"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
	Jti      string `json:"jti"`
}

var signingKey []byte

// SetKey sets the HMAC signing key. Call once at startup.
func SetKey(key string) {
	signingKey = []byte(key)
}

// Sign creates a signed JWT token.
func Sign(claims Claims) (string, error) {
	if claims.Jti == "" {
		claims.Jti = uuid.NewString()
	}
	if claims.Iat == 0 {
		claims.Iat = time.Now().Unix()
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signingInput := header + "." + payload
	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig, nil
}

// Verify validates a token and returns its claims.
func Verify(token string) (*Claims, error) {
	if len(signingKey) == 0 {
		return nil, fmt.Errorf("signing key not configured")
	}

	parts := splitToken(token)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Verify signature.
	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(signingInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload.
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	// Check expiry.
	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// NewClaims creates claims for a user with 24h expiry.
func NewClaims(userID int64, username, role string) Claims {
	now := time.Now()
	return Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Exp:      now.Add(24 * time.Hour).Unix(),
		Iat:      now.Unix(),
		Jti:      uuid.NewString(),
	}
}

func splitToken(token string) []string {
	parts := make([]string, 0, 3)
	start := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	if start < len(token) {
		parts = append(parts, token[start:])
	}
	return parts
}