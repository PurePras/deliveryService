package model

import "time"

type DeliverySlot struct {
	ID        string    `db:"id" json:"id"`
	Label     string    `db:"label" json:"label"`
	StartTime string    `db:"start_time" json:"start_time"` // TIME as "HH:MM:SS"
	EndTime   string    `db:"end_time" json:"end_time"`     // TIME as "HH:MM:SS"
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
