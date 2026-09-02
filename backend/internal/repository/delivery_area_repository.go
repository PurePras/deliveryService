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

type DeliveryAreaRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryAreaRepository(pool *pgxpool.Pool) *DeliveryAreaRepository {
	return &DeliveryAreaRepository{pool: pool}
}

const deliveryAreaColumns = "id, name, city, pincode, is_active, created_at, updated_at"

func scanDeliveryArea(row pgx.Row) (*model.DeliveryArea, error) {
	var a model.DeliveryArea
	err := row.Scan(&a.ID, &a.Name, &a.City, &a.Pincode, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *DeliveryAreaRepository) Create(ctx context.Context, a *model.DeliveryArea) (*model.DeliveryArea, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO delivery_areas (name, city, pincode, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING `+deliveryAreaColumns, a.Name, a.City, a.Pincode, a.IsActive)

	created, err := scanDeliveryArea(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: this area already exists for that pincode", apperror.ErrConflict)
		}
		return nil, err
	}
	return created, nil
}

func (r *DeliveryAreaRepository) GetByID(ctx context.Context, id string) (*model.DeliveryArea, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+deliveryAreaColumns+` FROM delivery_areas WHERE id = $1`, id)
	return scanDeliveryArea(row)
}

func (r *DeliveryAreaRepository) List(ctx context.Context, isActive *bool, limit, offset int) ([]model.DeliveryArea, error) {
	query := `SELECT ` + deliveryAreaColumns + ` FROM delivery_areas`
	var args []any
	if isActive != nil {
		args = append(args, *isActive)
		query += fmt.Sprintf(" WHERE is_active = $%d", len(args))
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY city, name LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	areas := []model.DeliveryArea{}
	for rows.Next() {
		var a model.DeliveryArea
		if err := rows.Scan(&a.ID, &a.Name, &a.City, &a.Pincode, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		areas = append(areas, a)
	}
	return areas, rows.Err()
}

func (r *DeliveryAreaRepository) Update(ctx context.Context, id string, a *model.DeliveryArea) (*model.DeliveryArea, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE delivery_areas
		SET name = $1, city = $2, pincode = $3, is_active = $4
		WHERE id = $5
		RETURNING `+deliveryAreaColumns, a.Name, a.City, a.Pincode, a.IsActive, id)

	updated, err := scanDeliveryArea(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: this area already exists for that pincode", apperror.ErrConflict)
		}
		return nil, err
	}
	return updated, nil
}

func (r *DeliveryAreaRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM delivery_areas WHERE id = $1`, id)
	if err != nil {
		if apperror.IsForeignKeyViolation(err) {
			return fmt.Errorf("%w: delivery area is still referenced by existing orders", apperror.ErrConflict)
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
