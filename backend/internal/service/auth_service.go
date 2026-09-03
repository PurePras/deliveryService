package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/auth"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
)

type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
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
