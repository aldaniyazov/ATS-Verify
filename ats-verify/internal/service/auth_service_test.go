package service

import (
	"testing"
	"time"

	"ats-verify/internal/middleware"
	"ats-verify/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "my-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash string")
	}

	// Verify the hash matches the password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Fatalf("expected password to match hash, got error: %v", err)
	}

	// Verify wrong password fails
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("wrong-password"))
	if err == nil {
		t.Fatal("expected error, got none for wrong password")
	}
}

// Minimal test to verify JWT token generation given a valid claims structure
func TestJWTGeneration(t *testing.T) {
	secret := "test-secret"
	claims := &middleware.Claims{
		UserID:            "123",
		Username:          "admin",
		Role:              models.RoleAdmin,
		MarketplacePrefix: "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}
}
