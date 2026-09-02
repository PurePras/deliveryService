package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
)

type DeliverySlotRepository struct {
	pool *pgxpool.Pool
}

func NewDeliverySlotRepository(pool *pgxpool.Pool) *DeliverySlotRepository {
	return &DeliverySlotRepository{pool: pool}
}

const deliverySlotColumns = "id, label, start_time::text, end_time::text, is_active, created_at, updated_at"

func scanDeliverySlot(row pgx.Row) (*model.DeliverySlot, error) {
	var s model.DeliverySlot
	err := row.Scan(&s.ID, &s.Label, &s.StartTime, &s.EndTime, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *DeliverySlotRepository) Create(ctx context.Context, s *model.DeliverySlot) (*model.DeliverySlot, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO delivery_slots (label, start_time, end_time, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING `+deliverySlotColumns, s.Label, s.StartTime, s.EndTime, s.IsActive)

	created, err := scanDeliverySlot(row)
	if err != nil {
		if apperror.IsCheckViolation(err) {
			return nil, fmt.Errorf("%w: end_time must be after start_time", apperror.ErrValidation)
		}
		return nil, err
	}
	return created, nil
}

func (r *DeliverySlotRepository) GetByID(ctx context.Context, id string) (*model.DeliverySlot, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+deliverySlotColumns+` FROM delivery_slots WHERE id = $1`, id)
	return scanDeliverySlot(row)
}

func (r *DeliverySlotRepository) List(ctx context.Context, isActive *bool, limit, offset int) ([]model.DeliverySlot, error) {
	query := `SELECT ` + deliverySlotColumns + ` FROM delivery_slots`
	var args []any
	if isActive != nil {
		args = append(args, *isActive)
		query += fmt.Sprintf(" WHERE is_active = $%d", len(args))
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY start_time LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := []model.DeliverySlot{}
	for rows.Next() {
		var s model.DeliverySlot
		if err := rows.Scan(&s.ID, &s.Label, &s.StartTime, &s.EndTime, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}

func (r *DeliverySlotRepository) Update(ctx context.Context, id string, s *model.DeliverySlot) (*model.DeliverySlot, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE delivery_slots
		SET label = $1, start_time = $2, end_time = $3, is_active = $4
		WHERE id = $5
		RETURNING `+deliverySlotColumns, s.Label, s.StartTime, s.EndTime, s.IsActive, id)

	updated, err := scanDeliverySlot(row)
	if err != nil {
		if apperror.IsCheckViolation(err) {
			return nil, fmt.Errorf("%w: end_time must be after start_time", apperror.ErrValidation)
		}
		return nil, err
	}
	return updated, nil
}

func (r *DeliverySlotRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM delivery_slots WHERE id = $1`, id)
	if err != nil {
		if apperror.IsForeignKeyViolation(err) {
			return fmt.Errorf("%w: delivery slot is still referenced by existing orders", apperror.ErrConflict)
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
