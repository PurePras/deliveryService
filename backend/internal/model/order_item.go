package model

import "time"

type OrderItem struct {
	ID        string    `db:"id" json:"id"`
	OrderID   string    `db:"order_id" json:"order_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Quantity  string    `db:"quantity" json:"quantity"`     // NUMERIC as string; fractional units (e.g. 0.5 kg) allowed
	UnitPrice string    `db:"unit_price" json:"unit_price"` // price snapshot at order time
	Subtotal  string    `db:"subtotal" json:"subtotal"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
