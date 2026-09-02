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

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

const categoryColumns = "id, name, slug, description, image_url, is_active, created_at, updated_at"

func scanCategory(row pgx.Row) (*model.Category, error) {
	var c model.Category
	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ImageURL, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) (*model.Category, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO categories (name, slug, description, image_url, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+categoryColumns, c.Name, c.Slug, c.Description, c.ImageURL, c.IsActive)

	created, err := scanCategory(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: a category with this name or slug already exists", apperror.ErrConflict)
		}
		return nil, err
	}
	return created, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*model.Category, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+categoryColumns+` FROM categories WHERE id = $1`, id)
	return scanCategory(row)
}

func (r *CategoryRepository) List(ctx context.Context, isActive *bool, limit, offset int) ([]model.Category, error) {
	query := `SELECT ` + categoryColumns + ` FROM categories`
	var args []any
	if isActive != nil {
		args = append(args, *isActive)
		query += fmt.Sprintf(" WHERE is_active = $%d", len(args))
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY name LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []model.Category{}
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ImageURL, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, id string, c *model.Category) (*model.Category, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE categories
		SET name = $1, slug = $2, description = $3, image_url = $4, is_active = $5
		WHERE id = $6
		RETURNING `+categoryColumns, c.Name, c.Slug, c.Description, c.ImageURL, c.IsActive, id)

	updated, err := scanCategory(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: a category with this name or slug already exists", apperror.ErrConflict)
		}
		return nil, err
	}
	return updated, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		if apperror.IsForeignKeyViolation(err) {
			return fmt.Errorf("%w: category is still referenced by existing products", apperror.ErrConflict)
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
