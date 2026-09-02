package service

import (
	"context"

	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Get(ctx context.Context, id string) (*model.User, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}
