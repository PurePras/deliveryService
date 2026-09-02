package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
)

var validProductUnits = map[string]bool{"kg": true, "g": true, "piece": true, "bundle": true, "dozen": true}

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, in dto.ProductInput) (*model.Product, error) {
	if err := validateProductInput(in); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, productFromInput(in))
}

func (s *ProductService) Get(ctx context.Context, id string) (*model.Product, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, filter repository.ProductFilter, limit, offset int) ([]model.Product, error) {
	if filter.CategoryID != nil {
		if err := requireUUID("category_id", *filter.CategoryID); err != nil {
			return nil, err
		}
	}
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *ProductService) Update(ctx context.Context, id string, in dto.ProductInput) (*model.Product, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	if err := validateProductInput(in); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, productFromInput(in))
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	if err := requireUUID("id", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func productFromInput(in dto.ProductInput) *model.Product {
	isAvailable := true
	if in.IsAvailable != nil {
		isAvailable = *in.IsAvailable
	}
	return &model.Product{
		CategoryID:    in.CategoryID,
		Name:          in.Name,
		Slug:          in.Slug,
		Description:   in.Description,
		Unit:          in.Unit,
		Price:         in.Price,
		StockQuantity: in.StockQuantity,
		ImageURL:      in.ImageURL,
		IsAvailable:   isAvailable,
	}
}

func validateProductInput(in dto.ProductInput) error {
	errs := []error{
		requireUUID("category_id", in.CategoryID),
		requireNonEmpty("name", in.Name),
		requireNonEmpty("slug", in.Slug),
		requireNonNegativeNumber("price", in.Price),
		requireNonNegativeNumber("stock_quantity", in.StockQuantity),
	}
	if !validProductUnits[in.Unit] {
		errs = append(errs, fmt.Errorf("%w: unit must be one of kg, g, piece, bundle, dozen", apperror.ErrValidation))
	}
	return errors.Join(errs...)
}
