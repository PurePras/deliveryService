package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
)

type DeliverySlotService struct {
	repo *repository.DeliverySlotRepository
}

func NewDeliverySlotService(repo *repository.DeliverySlotRepository) *DeliverySlotService {
	return &DeliverySlotService{repo: repo}
}

func (s *DeliverySlotService) Create(ctx context.Context, in dto.DeliverySlotInput) (*model.DeliverySlot, error) {
	if err := validateDeliverySlotInput(in); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, deliverySlotFromInput(in))
}

func (s *DeliverySlotService) Get(ctx context.Context, id string) (*model.DeliverySlot, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DeliverySlotService) List(ctx context.Context, isActive *bool, limit, offset int) ([]model.DeliverySlot, error) {
	return s.repo.List(ctx, isActive, limit, offset)
}

func (s *DeliverySlotService) Update(ctx context.Context, id string, in dto.DeliverySlotInput) (*model.DeliverySlot, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	if err := validateDeliverySlotInput(in); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, deliverySlotFromInput(in))
}

func (s *DeliverySlotService) Delete(ctx context.Context, id string) error {
	if err := requireUUID("id", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func deliverySlotFromInput(in dto.DeliverySlotInput) *model.DeliverySlot {
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}
	return &model.DeliverySlot{
		Label:     in.Label,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		IsActive:  isActive,
	}
}

func validateDeliverySlotInput(in dto.DeliverySlotInput) error {
	return errors.Join(
		requireNonEmpty("label", in.Label),
		requireTimeOfDay("start_time", in.StartTime),
		requireTimeOfDay("end_time", in.EndTime),
	)
}

func requireTimeOfDay(field, value string) error {
	for _, layout := range []string{"15:04", "15:04:05"} {
		if _, err := time.Parse(layout, value); err == nil {
			return nil
		}
	}
	return fmt.Errorf("%w: %s must be a time in HH:MM or HH:MM:SS format", apperror.ErrValidation, field)
}
