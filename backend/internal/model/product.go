package model

import "time"

type Product struct {
	ID            string    `db:"id" json:"id"`
	CategoryID    string    `db:"category_id" json:"category_id"`
	Name          string    `db:"name" json:"name"`
	Slug          string    `db:"slug" json:"slug"`
	Description   *string   `db:"description" json:"description,omitempty"`
	Unit          string    `db:"unit" json:"unit"`
	Price         string    `db:"price" json:"price"`                   // NUMERIC as string to avoid float rounding on currency
	StockQuantity string    `db:"stock_quantity" json:"stock_quantity"` // NUMERIC as string; fractional units (e.g. 0.5 kg) allowed
	ImageURL      *string   `db:"image_url" json:"image_url,omitempty"`
	IsAvailable   bool      `db:"is_available" json:"is_available"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
