package password

import (
	"testing"
)

func TestHashAndCompare(t *testing.T) {
	plain := "SecretP@ss123"

	hash, err := Hash(plain)
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if err := Compare(hash, plain); err != nil {
		t.Fatalf("expected password match, got error: %v", err)
	}

	if err := Compare(hash, "wrong-password"); err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}
