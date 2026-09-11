package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/auth"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
	"github.com/PurePras/shri-ram-service/backend/internal/sms"
)

// OTP tuning: resend/attempt caps are fixed policy (like auth.TTL), while how long a
// code stays valid is the one thing worth making operator-configurable (Config.OTPTTL).
const (
	otpResendCooldown = 30 * time.Second
	otpMaxAttempts    = 5
	otpMaxPerHour     = 5
)

type AuthService struct {
	repo      *repository.UserRepository
	otpRepo   *repository.OTPRepository
	sms       sms.Sender
	jwtSecret string
	otpTTL    time.Duration
}

func NewAuthService(repo *repository.UserRepository, otpRepo *repository.OTPRepository, sender sms.Sender, jwtSecret string, otpTTL time.Duration) *AuthService {
	return &AuthService{repo: repo, otpRepo: otpRepo, sms: sender, jwtSecret: jwtSecret, otpTTL: otpTTL}
}

func (s *AuthService) Register(ctx context.Context, in dto.RegisterRequest) (*model.User, string, error) {
	errs := []error{
		requireNonEmpty("name", in.Name),
		requirePhone("phone", in.Phone),
		requireMinLength("password", in.Password, 8),
	}
	if in.Email != nil && *in.Email != "" {
		errs = append(errs, requireEmail("email", *in.Email))
	}
	if err := errors.Join(errs...); err != nil {
		return nil, "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	hashStr := string(hash)

	user, err := s.repo.Create(ctx, &model.User{
		Name:         in.Name,
		Phone:        in.Phone,
		Email:        in.Email,
		PasswordHash: &hashStr,
	})
	if err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login returns the same generic "invalid phone or password" error whether the
// phone doesn't exist or the password is wrong, so the response never reveals
// which one was incorrect.
func (s *AuthService) Login(ctx context.Context, in dto.LoginRequest) (*model.User, string, error) {
	if err := errors.Join(
		requireNonEmpty("phone", in.Phone),
		requireNonEmpty("password", in.Password),
	); err != nil {
		return nil, "", err
	}

	invalidCreds := fmt.Errorf("%w: invalid phone or password", apperror.ErrUnauthorized)

	user, err := s.repo.GetByPhone(ctx, in.Phone)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, "", invalidCreds
		}
		return nil, "", err
	}

	if user.PasswordHash == nil {
		return nil, "", invalidCreds
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, "", invalidCreds
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// ChangePassword re-verifies currentPassword the same way Login does before saving newPassword,
// so a stolen/left-open session alone isn't enough to take over the password.
func (s *AuthService) ChangePassword(ctx context.Context, userID string, in dto.ChangePasswordRequest) error {
	if err := errors.Join(
		requireNonEmpty("current_password", in.CurrentPassword),
		requireMinLength("new_password", in.NewPassword, 8),
	); err != nil {
		return err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.PasswordHash == nil {
		return fmt.Errorf("%w: current password is incorrect", apperror.ErrUnauthorized)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(in.CurrentPassword)); err != nil {
		return fmt.Errorf("%w: current password is incorrect", apperror.ErrUnauthorized)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePasswordHash(ctx, userID, string(hash))
}

// RequestLoginOTP validates the phone and, subject to the resend cooldown and
// hourly cap, generates and sends a fresh code. Like Login's generic "invalid
// phone or password", this never reveals whether the phone has an account —
// it returns success either way and simply sends nothing for an unregistered
// number, so the endpoint can't be used to enumerate accounts.
func (s *AuthService) RequestLoginOTP(ctx context.Context, phone string) error {
	if err := requirePhone("phone", phone); err != nil {
		return err
	}

	_, err := s.repo.GetByPhone(ctx, phone)
	switch {
	case err == nil:
		// registered — proceed to send below.
	case errors.Is(err, apperror.ErrNotFound):
		// unregistered — still enforce the same limits (see below) before
		// no-op'ing, so this can't be used to SMS-bomb an arbitrary number either.
	default:
		return err
	}
	registered := err == nil

	latest, err := s.otpRepo.GetLatestByPhone(ctx, phone)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return err
	}
	if latest != nil {
		if wait := otpResendCooldown - time.Since(latest.CreatedAt); wait > 0 {
			return fmt.Errorf("%w: please wait %d seconds before requesting another OTP",
				apperror.ErrValidation, int(wait.Seconds()+1))
		}
	}

	count, err := s.otpRepo.CountSince(ctx, phone, time.Now().Add(-time.Hour))
	if err != nil {
		return err
	}
	if count >= otpMaxPerHour {
		return fmt.Errorf("%w: too many OTP requests for this number, try again later", apperror.ErrValidation)
	}

	if !registered {
		return nil
	}

	code, err := generateOTPCode()
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, err := s.otpRepo.Create(ctx, phone, string(hash), time.Now().Add(s.otpTTL)); err != nil {
		return err
	}

	return s.sms.Send(ctx, phone, code)
}

// VerifyLoginOTP checks code against the most recently generated OTP for phone —
// requesting a new one supersedes whatever came before, so only the latest is ever
// valid — and, if it matches, issues a session token exactly like Login does.
func (s *AuthService) VerifyLoginOTP(ctx context.Context, phone, code string) (*model.User, string, error) {
	if err := errors.Join(
		requirePhone("phone", phone),
		requireOTPCode("code", code),
	); err != nil {
		return nil, "", err
	}

	invalidOTP := fmt.Errorf("%w: invalid or expired code", apperror.ErrUnauthorized)

	latest, err := s.otpRepo.GetLatestByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, "", invalidOTP
		}
		return nil, "", err
	}

	if latest.ConsumedAt != nil || time.Now().After(latest.ExpiresAt) {
		return nil, "", invalidOTP
	}
	if latest.Attempts >= otpMaxAttempts {
		return nil, "", fmt.Errorf("%w: too many incorrect attempts, request a new code", apperror.ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(latest.CodeHash), []byte(code)); err != nil {
		if err := s.otpRepo.IncrementAttempts(ctx, latest.ID); err != nil {
			return nil, "", err
		}
		return nil, "", invalidOTP
	}

	if err := s.otpRepo.MarkConsumed(ctx, latest.ID); err != nil {
		return nil, "", err
	}

	user, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, "", invalidOTP
		}
		return nil, "", err
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// generateOTPCode uses crypto/rand, not math/rand — this is a security credential,
// not test fixture data.
func generateOTPCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
