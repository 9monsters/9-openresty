package auth

import (
	"testing"
	"time"
)

func TestSignVerifyRoundTrip(t *testing.T) {
	SetKey("test-signing-key")
	claims := NewClaims(42, "admin", "Administrator")
	tok, err := Sign(claims)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	got, err := Verify(tok)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got.UserID != 42 || got.Username != "admin" || got.Role != "Administrator" {
		t.Fatalf("claims mismatch: %+v", got)
	}
}

func TestVerifyTamperedSignature(t *testing.T) {
	SetKey("test-signing-key")
	tok, _ := Sign(NewClaims(1, "admin", "Administrator"))
	tampered := tok[:len(tok)-1] + "x"
	if _, err := Verify(tampered); err == nil {
		t.Fatal("expected error for a tampered signature")
	}
}

func TestVerifyExpiredToken(t *testing.T) {
	SetKey("test-signing-key")
	c := NewClaims(1, "admin", "Administrator")
	c.Exp = time.Now().Add(-time.Hour).Unix()
	tok, _ := Sign(c)
	if _, err := Verify(tok); err == nil {
		t.Fatal("expected error for an expired token")
	}
}

func TestVerifyWrongKey(t *testing.T) {
	SetKey("key-a")
	tok, _ := Sign(NewClaims(1, "admin", "Administrator"))
	SetKey("key-b")
	if _, err := Verify(tok); err == nil {
		t.Fatal("expected error when verifying with a different key")
	}
}

func TestNewClaimsSetsExpiryAndJti(t *testing.T) {
	c := NewClaims(7, "op", "Operator")
	if c.Exp <= time.Now().Unix() {
		t.Fatal("expected expiry in the future")
	}
	if c.UserID != 7 || c.Username != "op" || c.Role != "Operator" {
		t.Fatalf("unexpected claims: %+v", c)
	}
	if c.Jti == "" {
		t.Fatal("expected a jti to be assigned")
	}
}
