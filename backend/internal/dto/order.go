package dto

import "github.com/PurePras/shri-ram-service/backend/internal/model"

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  string `json:"quantity"`
}

type CreateOrderRequest struct {
	DeliveryAreaID  string                   `json:"delivery_area_id"`
	DeliverySlotID  string                   `json:"delivery_slot_id"`
	DeliveryDate    string                   `json:"delivery_date"` // "YYYY-MM-DD"
	DeliveryAddress string                   `json:"delivery_address"`
	Items           []CreateOrderItemRequest `json:"items"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type OrderResponse struct {
	model.Order
	Items []model.OrderItem `json:"items"`
}
