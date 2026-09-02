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

var validOrderStatuses = map[string]bool{
	model.OrderStatusPending:        true,
	model.OrderStatusConfirmed:      true,
	model.OrderStatusOutForDelivery: true,
	model.OrderStatusDelivered:      true,
	model.OrderStatusCancelled:      true,
}

var terminalOrderStatuses = map[string]bool{
	model.OrderStatusDelivered: true,
	model.OrderStatusCancelled: true,
}

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, in dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	deliveryDate, err := validateCreateOrderRequest(in)
	if err != nil {
		return nil, err
	}

	items := make([]repository.OrderItemInput, len(in.Items))
	for i, item := range in.Items {
		items[i] = repository.OrderItemInput{ProductID: item.ProductID, Quantity: item.Quantity}
	}

	order, orderItems, err := s.repo.Create(ctx, &model.Order{
		UserID:          in.UserID,
		DeliveryAreaID:  in.DeliveryAreaID,
		DeliverySlotID:  in.DeliverySlotID,
		DeliveryDate:    deliveryDate,
		DeliveryAddress: in.DeliveryAddress,
	}, items)
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{Order: *order, Items: orderItems}, nil
}

func (s *OrderService) Get(ctx context.Context, id string) (*dto.OrderResponse, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	order, items, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.OrderResponse{Order: *order, Items: items}, nil
}

func (s *OrderService) List(ctx context.Context, userID *string, limit, offset int) ([]model.Order, error) {
	if userID != nil {
		if err := requireUUID("user_id", *userID); err != nil {
			return nil, err
		}
	}
	return s.repo.List(ctx, userID, limit, offset)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id string, in dto.UpdateOrderStatusRequest) (*model.Order, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	if !validOrderStatuses[in.Status] {
		return nil, fmt.Errorf("%w: status must be one of pending, confirmed, out_for_delivery, delivered, cancelled", apperror.ErrValidation)
	}

	current, _, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if terminalOrderStatuses[current.Status] {
		return nil, fmt.Errorf("%w: order is already %s and cannot change status", apperror.ErrConflict, current.Status)
	}

	return s.repo.UpdateStatus(ctx, id, in.Status)
}

func validateCreateOrderRequest(in dto.CreateOrderRequest) (time.Time, error) {
	errs := []error{
		requireUUID("user_id", in.UserID),
		requireUUID("delivery_area_id", in.DeliveryAreaID),
		requireUUID("delivery_slot_id", in.DeliverySlotID),
		requireNonEmpty("delivery_address", in.DeliveryAddress),
	}

	deliveryDate, dateErr := time.Parse("2006-01-02", in.DeliveryDate)
	if dateErr != nil {
		errs = append(errs, fmt.Errorf("%w: delivery_date must be in YYYY-MM-DD format", apperror.ErrValidation))
	}

	if len(in.Items) == 0 {
		errs = append(errs, fmt.Errorf("%w: at least one item is required", apperror.ErrValidation))
	}
	for _, item := range in.Items {
		errs = append(errs,
			requireUUID("items.product_id", item.ProductID),
			requirePositiveNumber("items.quantity", item.Quantity),
		)
	}

	return deliveryDate, errors.Join(errs...)
}
