package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("s3cret!")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if hash == "s3cret!" {
		t.Fatal("hash must not equal the plaintext password")
	}
	if !CheckPassword(hash, "s3cret!") {
		t.Fatal("expected CheckPassword to accept the correct password")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("expected CheckPassword to reject a wrong password")
	}
}

func TestHashPasswordProducesUniqueSalts(t *testing.T) {
	h1, err1 := HashPassword("same-password")
	h2, err2 := HashPassword("same-password")
	if err1 != nil || err2 != nil {
		t.Fatalf("hash errors: %v, %v", err1, err2)
	}
	if h1 == h2 {
		t.Fatal("bcrypt should use a unique salt per hash")
	}
}

func TestCheckPasswordEmptyHash(t *testing.T) {
	if CheckPassword("", "anything") {
		t.Fatal("expected empty hash to reject all passwords")
	}
}
