//go:build integration

package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

const testJWTSecret = "integration-test-jwt-secret"

func TestAuthFlow_RegisterLoginChangePassword(t *testing.T) {
	ctx := context.Background()
	svc := service.NewAuthService(repository.NewUserRepository(testPool), testJWTSecret)

	phone := randomTestPhone()
	const initialPassword = "correct-password-123"
	const newPassword = "new-password-456"

	registered, token, err := svc.Register(ctx, dto.RegisterRequest{
		Name:     "Integration Test User",
		Phone:    phone,
		Password: initialPassword,
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token from Register")
	}
	t.Cleanup(func() { cleanupUser(t, registered.ID) })

	loggedIn, _, err := svc.Login(ctx, dto.LoginRequest{Phone: phone, Password: initialPassword})
	if err != nil {
		t.Fatalf("Login with the just-registered password: %v", err)
	}
	if loggedIn.ID != registered.ID {
		t.Fatalf("Login returned user %s, want %s", loggedIn.ID, registered.ID)
	}

	if _, _, err := svc.Login(ctx, dto.LoginRequest{Phone: phone, Password: "wrong-password"}); err == nil {
		t.Fatal("expected an error logging in with the wrong password")
	} else if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}

	if err := svc.ChangePassword(ctx, registered.ID, dto.ChangePasswordRequest{
		CurrentPassword: initialPassword,
		NewPassword:     newPassword,
	}); err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}

	if _, _, err := svc.Login(ctx, dto.LoginRequest{Phone: phone, Password: initialPassword}); err == nil {
		t.Fatal("expected the old password to be rejected after ChangePassword")
	}

	if _, _, err := svc.Login(ctx, dto.LoginRequest{Phone: phone, Password: newPassword}); err != nil {
		t.Fatalf("Login with the new password after ChangePassword: %v", err)
	}
}

func TestAuthFlow_DuplicatePhoneRejected(t *testing.T) {
	ctx := context.Background()
	svc := service.NewAuthService(repository.NewUserRepository(testPool), testJWTSecret)

	phone := randomTestPhone()
	first, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "First User", Phone: phone, Password: "correct-password-123",
	})
	if err != nil {
		t.Fatalf("Register (first): %v", err)
	}
	t.Cleanup(func() { cleanupUser(t, first.ID) })

	if _, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "Second User", Phone: phone, Password: "another-password-456",
	}); err == nil {
		t.Fatal("expected an error registering a second user with an already-taken phone number")
	} else if !errors.Is(err, apperror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got: %v", err)
	}
}

func cleanupUser(t *testing.T, id string) {
	t.Helper()
	if _, err := testPool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id); err != nil {
		t.Logf("cleanup: delete user %s: %v", id, err)
	}
}
