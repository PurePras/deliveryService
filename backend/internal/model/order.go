package model

import "time"

const (
	OrderStatusPending        = "pending"
	OrderStatusConfirmed      = "confirmed"
	OrderStatusOutForDelivery = "out_for_delivery"
	OrderStatusDelivered      = "delivered"
	OrderStatusCancelled      = "cancelled"
)

type Order struct {
	ID              string    `db:"id" json:"id"`
	UserID          string    `db:"user_id" json:"user_id"`
	DeliveryAreaID  string    `db:"delivery_area_id" json:"delivery_area_id"`
	DeliverySlotID  string    `db:"delivery_slot_id" json:"delivery_slot_id"`
	DeliveryDate    time.Time `db:"delivery_date" json:"delivery_date"`
	DeliveryAddress string    `db:"delivery_address" json:"delivery_address"`
	Status          string    `db:"status" json:"status"`
	TotalAmount     string    `db:"total_amount" json:"total_amount"` // NUMERIC as string to avoid float rounding on currency
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
