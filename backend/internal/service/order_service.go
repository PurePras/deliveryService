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

// Create places an order for userID, which must come from the authenticated
// session (see middleware.UserFromContext) — never from client-supplied input.
func (s *OrderService) Create(ctx context.Context, userID string, in dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if err := requireUUID("user_id", userID); err != nil {
		return nil, err
	}
	deliveryDate, err := validateCreateOrderRequest(in)
	if err != nil {
		return nil, err
	}

	items := make([]repository.OrderItemInput, len(in.Items))
	for i, item := range in.Items {
		items[i] = repository.OrderItemInput{ProductID: item.ProductID, Quantity: item.Quantity}
	}

	order, orderItems, err := s.repo.Create(ctx, &model.Order{
		UserID:          userID,
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

// Get returns the order only if requesterID owns it or isAdmin is true; otherwise it
// returns apperror.ErrNotFound (never a 403) so a non-owned id can't be distinguished
// from one that simply doesn't exist.
func (s *OrderService) Get(ctx context.Context, id, requesterID string, isAdmin bool) (*dto.OrderResponse, error) {
	if err := requireUUID("id", id); err != nil {
		return nil, err
	}
	order, items, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order.UserID != requesterID && !isAdmin {
		return nil, apperror.ErrNotFound
	}
	return &dto.OrderResponse{Order: *order, Items: items}, nil
}

type OrderListParams struct {
	Status *string
	Limit  int
	Offset int
}

// List returns requesterID's own orders, unless isAdmin — in which case it returns every
// customer's orders (optionally filtered by status), since fulfillment needs to see all of them.
func (s *OrderService) List(ctx context.Context, requesterID string, isAdmin bool, params OrderListParams) ([]model.Order, error) {
	if err := requireUUID("user_id", requesterID); err != nil {
		return nil, err
	}
	filter := repository.OrderFilter{Status: params.Status}
	if !isAdmin {
		filter.UserID = &requesterID
	}
	return s.repo.List(ctx, filter, params.Limit, params.Offset)
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
		requireUUID("delivery_area_id", in.DeliveryAreaID),
		requireUUID("delivery_slot_id", in.DeliverySlotID),
		requireNonEmpty("delivery_address", in.DeliveryAddress),
	}

	deliveryDate, dateErr := time.Parse("2006-01-02", in.DeliveryDate)
	if dateErr != nil {
		errs = append(errs, fmt.Errorf("%w: delivery_date must be in YYYY-MM-DD format", apperror.ErrValidation))
	} else if today := time.Now().UTC(); deliveryDate.Before(time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)) {
		// time.Parse of a bare "YYYY-MM-DD" layout yields a UTC-midnight time.Time,
		// so "today" must be computed in UTC too to avoid a timezone-boundary off-by-one.
		errs = append(errs, fmt.Errorf("%w: delivery_date cannot be in the past", apperror.ErrValidation))
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
