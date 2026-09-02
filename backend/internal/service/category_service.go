package service

import (
	"context"
	"errors"

	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, in dto.CategoryInput) (*model.Category, error) {
	if err := validateCategoryInput(in); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, categoryFromInput(in))
}

func (s *CategoryService) Get(ctx context.Context, id string) (*model.Category, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, isActive *bool, limit, offset int) ([]model.Category, error) {
	return s.repo.List(ctx, isActive, limit, offset)
}

func (s *CategoryService) Update(ctx context.Context, id string, in dto.CategoryInput) (*model.Category, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	if err := validateCategoryInput(in); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, categoryFromInput(in))
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	if err := requireUUID("id", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func categoryFromInput(in dto.CategoryInput) *model.Category {
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}
	return &model.Category{
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		ImageURL:    in.ImageURL,
		IsActive:    isActive,
	}
}

func validateCategoryInput(in dto.CategoryInput) error {
	return errors.Join(
		requireNonEmpty("name", in.Name),
		requireNonEmpty("slug", in.Slug),
	)
}
