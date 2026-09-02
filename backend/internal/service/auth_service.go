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
	if err := errors.Join(
		requireNonEmpty("name", in.Name),
		requireNonEmpty("phone", in.Phone),
		requireMinLength("password", in.Password, 8),
	); err != nil {
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
