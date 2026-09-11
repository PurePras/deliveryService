package service

import (
	"errors"
	"testing"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
)

func TestRequireUUID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid v4", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid v4 without hyphens", "550e8400e29b41d4a716446655440000", false},
		{"empty", "", true},
		{"not a uuid", "not-a-uuid", true},
		{"too short", "550e8400-e29b-41d4-a716", true},
		{"invalid hex characters", "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireUUID("id", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequireNonEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"non-empty", "hello", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireNonEmpty("field", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequireNonNegativeNumber(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"positive", "10.50", false},
		{"zero", "0", false},
		{"negative", "-5", true},
		{"not a number", "abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireNonNegativeNumber("field", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequirePositiveNumber(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"positive", "0.5", false},
		{"zero", "0", true},
		{"negative", "-1", true},
		{"not a number", "abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requirePositiveNumber("field", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequireMinLength(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		min     int
		wantErr bool
	}{
		{"exact length", "12345678", 8, false},
		{"longer", "123456789", 8, false},
		{"too short", "1234567", 8, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireMinLength("field", tt.value, tt.min)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequirePhone(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid starting with 9", "9876500001", false},
		{"valid starting with 6", "6123456789", false},
		{"too short", "987650000", true},
		{"too long", "98765000011", true},
		{"starts with 5", "5876500001", true},
		{"starts with 0", "0876500001", true},
		{"contains letters", "98765abcde", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requirePhone("phone", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequireEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.co.in", false},
		{"missing @", "userexample.com", true},
		{"missing domain", "user@", true},
		{"missing tld", "user@example", true},
		{"contains space", "user name@example.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireEmail("email", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

func TestRequireOTPCode(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 6 digits", "123456", false},
		{"valid with leading zero", "012345", false},
		{"too short", "12345", true},
		{"too long", "1234567", true},
		{"contains letters", "12345a", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requireOTPCode("code", tt.value)
			assertValidationErr(t, err, tt.wantErr)
		})
	}
}

// assertValidationErr checks that err is nil when wantErr is false, and that it wraps
// apperror.ErrValidation when wantErr is true — every validator in this package must
// use that sentinel so handleServiceError can map it to a 400.
func assertValidationErr(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, apperror.ErrValidation) {
			t.Fatalf("expected error to wrap apperror.ErrValidation, got: %v", err)
		}
		return
	}
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
