package service

import (
	"context"
	"errors"

	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
)

type DeliveryAreaService struct {
	repo *repository.DeliveryAreaRepository
}

func NewDeliveryAreaService(repo *repository.DeliveryAreaRepository) *DeliveryAreaService {
	return &DeliveryAreaService{repo: repo}
}

func (s *DeliveryAreaService) Create(ctx context.Context, in dto.DeliveryAreaInput) (*model.DeliveryArea, error) {
	if err := validateDeliveryAreaInput(in); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, deliveryAreaFromInput(in))
}

func (s *DeliveryAreaService) Get(ctx context.Context, id string) (*model.DeliveryArea, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DeliveryAreaService) List(ctx context.Context, isActive *bool, limit, offset int) ([]model.DeliveryArea, error) {
	return s.repo.List(ctx, isActive, limit, offset)
}

func (s *DeliveryAreaService) Update(ctx context.Context, id string, in dto.DeliveryAreaInput) (*model.DeliveryArea, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	if err := validateDeliveryAreaInput(in); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, deliveryAreaFromInput(in))
}

func (s *DeliveryAreaService) Delete(ctx context.Context, id string) error {
	if err := requireUUID("id", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func deliveryAreaFromInput(in dto.DeliveryAreaInput) *model.DeliveryArea {
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}
	return &model.DeliveryArea{
		Name:     in.Name,
		City:     in.City,
		Pincode:  in.Pincode,
		IsActive: isActive,
	}
}

func validateDeliveryAreaInput(in dto.DeliveryAreaInput) error {
	return errors.Join(
		requireNonEmpty("name", in.Name),
		requireNonEmpty("city", in.City),
		requireNonEmpty("pincode", in.Pincode),
	)
}
