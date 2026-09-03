package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-for-unit-tests-only"

func TestGenerateAndParseToken_RoundTrip(t *testing.T) {
	token, err := GenerateToken("user-123", "customer", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token")
	}

	claims, err := ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-123")
	}
	if claims.Role != "customer" {
		t.Errorf("Role = %q, want %q", claims.Role, "customer")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken("user-123", "customer", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	if _, err := ParseToken(token, "a-completely-different-secret"); err == nil {
		t.Fatal("expected an error parsing a token signed with a different secret")
	}
}

func TestParseToken_Tampered(t *testing.T) {
	token, err := GenerateToken("user-123", "customer", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	// Flip the first character of the signature segment (not the last — the final
	// base64url character(s) can carry unused padding bits, so mutating one of those
	// isn't guaranteed to change the decoded signature bytes).
	sigStart := strings.LastIndex(token, ".") + 1
	c := token[sigStart]
	replacement := byte('A')
	if c == replacement {
		replacement = 'B'
	}
	tampered := token[:sigStart] + string(replacement) + token[sigStart+1:]

	if _, err := ParseToken(tampered, testSecret); err == nil {
		t.Fatal("expected an error parsing a tampered token")
	}
}

func TestParseToken_Expired(t *testing.T) {
	claims := Claims{
		UserID: "user-123",
		Role:   "customer",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("signing test token: %v", err)
	}

	if _, err := ParseToken(signed, testSecret); err == nil {
		t.Fatal("expected an error parsing an expired token")
	}
}
