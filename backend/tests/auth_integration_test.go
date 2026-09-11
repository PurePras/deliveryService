//go:build integration

package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

const testJWTSecret = "integration-test-jwt-secret"

// capturingSender stands in for a real SMS provider — none is wired into this
// project yet (see internal/sms) — so tests can read the code AuthService generated
// instead of intercepting an actual text message.
type capturingSender struct {
	lastPhone string
	lastCode  string
}

func (s *capturingSender) Send(_ context.Context, phone, code string) error {
	s.lastPhone, s.lastCode = phone, code
	return nil
}

func newTestAuthService(t *testing.T, otpTTL time.Duration) (*service.AuthService, *capturingSender) {
	t.Helper()
	sender := &capturingSender{}
	svc := service.NewAuthService(
		repository.NewUserRepository(testPool),
		repository.NewOTPRepository(testPool),
		sender,
		testJWTSecret,
		otpTTL,
	)
	return svc, sender
}

func TestAuthFlow_RegisterLoginChangePassword(t *testing.T) {
	ctx := context.Background()
	svc, _ := newTestAuthService(t, 5*time.Minute)

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
	svc, _ := newTestAuthService(t, 5*time.Minute)

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

func TestAuthFlow_OTPLogin_Success(t *testing.T) {
	ctx := context.Background()
	svc, sender := newTestAuthService(t, 5*time.Minute)

	phone := randomTestPhone()
	registered, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "OTP Test User", Phone: phone, Password: "correct-password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	t.Cleanup(func() { cleanupUser(t, registered.ID) })

	if err := svc.RequestLoginOTP(ctx, phone); err != nil {
		t.Fatalf("RequestLoginOTP: %v", err)
	}
	if sender.lastPhone != phone || sender.lastCode == "" {
		t.Fatalf("expected an OTP to be sent to %s, sender saw phone=%q code=%q", phone, sender.lastPhone, sender.lastCode)
	}

	user, token, err := svc.VerifyLoginOTP(ctx, phone, sender.lastCode)
	if err != nil {
		t.Fatalf("VerifyLoginOTP with the correct code: %v", err)
	}
	if user.ID != registered.ID {
		t.Fatalf("VerifyLoginOTP returned user %s, want %s", user.ID, registered.ID)
	}
	if token == "" {
		t.Fatal("expected a non-empty token from VerifyLoginOTP")
	}

	// Single-use: the same code must not verify a second time.
	if _, _, err := svc.VerifyLoginOTP(ctx, phone, sender.lastCode); err == nil {
		t.Fatal("expected an error reusing an already-consumed OTP")
	} else if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestAuthFlow_OTPLogin_WrongCodeRejected(t *testing.T) {
	ctx := context.Background()
	svc, sender := newTestAuthService(t, 5*time.Minute)

	phone := randomTestPhone()
	registered, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "OTP Test User", Phone: phone, Password: "correct-password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	t.Cleanup(func() { cleanupUser(t, registered.ID) })

	if err := svc.RequestLoginOTP(ctx, phone); err != nil {
		t.Fatalf("RequestLoginOTP: %v", err)
	}

	wrongCode := "000000"
	if sender.lastCode == wrongCode {
		wrongCode = "111111"
	}
	if _, _, err := svc.VerifyLoginOTP(ctx, phone, wrongCode); err == nil {
		t.Fatal("expected an error verifying the wrong code")
	} else if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestAuthFlow_OTPLogin_ExcessiveAttemptsRejected(t *testing.T) {
	ctx := context.Background()
	svc, sender := newTestAuthService(t, 5*time.Minute)

	phone := randomTestPhone()
	registered, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "OTP Test User", Phone: phone, Password: "correct-password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	t.Cleanup(func() { cleanupUser(t, registered.ID) })

	if err := svc.RequestLoginOTP(ctx, phone); err != nil {
		t.Fatalf("RequestLoginOTP: %v", err)
	}
	wrongCode := "000000"
	if sender.lastCode == wrongCode {
		wrongCode = "111111"
	}

	// Exhaust the attempt cap with wrong codes...
	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		if _, _, err := svc.VerifyLoginOTP(ctx, phone, wrongCode); err == nil {
			t.Fatalf("attempt %d: expected an error verifying the wrong code", i+1)
		}
	}

	// ...then even the correct code must now be rejected as "too many attempts",
	// not silently accepted.
	if _, _, err := svc.VerifyLoginOTP(ctx, phone, sender.lastCode); err == nil {
		t.Fatal("expected the correct code to still be rejected once the attempt cap is hit")
	} else if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestAuthFlow_OTPLogin_Expired(t *testing.T) {
	ctx := context.Background()
	svc, sender := newTestAuthService(t, 1*time.Millisecond)

	phone := randomTestPhone()
	registered, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "OTP Test User", Phone: phone, Password: "correct-password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	t.Cleanup(func() { cleanupUser(t, registered.ID) })

	if err := svc.RequestLoginOTP(ctx, phone); err != nil {
		t.Fatalf("RequestLoginOTP: %v", err)
	}
	time.Sleep(10 * time.Millisecond)

	if _, _, err := svc.VerifyLoginOTP(ctx, phone, sender.lastCode); err == nil {
		t.Fatal("expected an error verifying an expired code")
	} else if !errors.Is(err, apperror.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestAuthFlow_OTPLogin_ResendCooldown(t *testing.T) {
	ctx := context.Background()
	svc, _ := newTestAuthService(t, 5*time.Minute)

	phone := randomTestPhone()
	registered, _, err := svc.Register(ctx, dto.RegisterRequest{
		Name: "OTP Test User", Phone: phone, Password: "correct-password-123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	t.Cleanup(func() { cleanupUser(t, registered.ID) })

	if err := svc.RequestLoginOTP(ctx, phone); err != nil {
		t.Fatalf("RequestLoginOTP (first): %v", err)
	}
	if err := svc.RequestLoginOTP(ctx, phone); err == nil {
		t.Fatal("expected an error requesting a second OTP within the resend cooldown")
	} else if !errors.Is(err, apperror.ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}
}

func TestAuthFlow_OTPLogin_UnregisteredPhoneSendsNothing(t *testing.T) {
	ctx := context.Background()
	svc, sender := newTestAuthService(t, 5*time.Minute)

	// Never registered — RequestLoginOTP must still report success (so the endpoint
	// can't be used to discover which numbers have accounts) but must not actually
	// send anything.
	phone := randomTestPhone()
	if err := svc.RequestLoginOTP(ctx, phone); err != nil {
		t.Fatalf("RequestLoginOTP for an unregistered phone: %v", err)
	}
	if sender.lastCode != "" {
		t.Fatal("expected no OTP to be sent for an unregistered phone")
	}
}

func cleanupUser(t *testing.T, id string) {
	t.Helper()
	if _, err := testPool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id); err != nil {
		t.Logf("cleanup: delete user %s: %v", id, err)
	}
}
