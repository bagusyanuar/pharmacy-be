package jwt

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndParse(t *testing.T) {
	manager := NewManager("test-secret-key-12345", 1*time.Hour, "test-issuer")

	userID := "user-uuid-12345"
	token, err := manager.GenerateToken(userID)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("expected no error parsing token, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	manager := NewManager("test-secret-key-12345", 1*time.Hour, "test-issuer")

	_, err := manager.ParseToken("invalid.token.here")
	if err == nil {
		t.Fatal("expected error on invalid token, got nil")
	}
}
