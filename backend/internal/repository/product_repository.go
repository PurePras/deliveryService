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

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

const productColumns = `id, category_id, name, slug, description, unit, price::text, stock_quantity::text,
	image_url, is_available, created_at, updated_at`

func scanProduct(row pgx.Row) (*model.Product, error) {
	var p model.Product
	err := row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Slug, &p.Description, &p.Unit, &p.Price, &p.StockQuantity,
		&p.ImageURL, &p.IsAvailable, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) (*model.Product, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO products (category_id, name, slug, description, unit, price, stock_quantity, image_url, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING `+productColumns,
		p.CategoryID, p.Name, p.Slug, p.Description, p.Unit, p.Price, p.StockQuantity, p.ImageURL, p.IsAvailable)

	created, err := scanProduct(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: a product with this slug already exists", apperror.ErrConflict)
		}
		if apperror.IsForeignKeyViolation(err) {
			return nil, fmt.Errorf("%w: category does not exist", apperror.ErrValidation)
		}
		if apperror.IsCheckViolation(err) {
			return nil, fmt.Errorf("%w: %v", apperror.ErrValidation, err)
		}
		return nil, err
	}
	return created, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+productColumns+` FROM products WHERE id = $1`, id)
	return scanProduct(row)
}

type ProductFilter struct {
	CategoryID  *string
	IsAvailable *bool
	Search      *string
}

func (r *ProductRepository) List(ctx context.Context, filter ProductFilter, limit, offset int) ([]model.Product, error) {
	query := `SELECT ` + productColumns + ` FROM products`
	var conditions []string
	var args []any
	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
	}
	if filter.IsAvailable != nil {
		args = append(args, *filter.IsAvailable)
		conditions = append(conditions, fmt.Sprintf("is_available = $%d", len(args)))
	}
	if filter.Search != nil && *filter.Search != "" {
		args = append(args, "%"+*filter.Search+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	for i, cond := range conditions {
		if i == 0 {
			query += " WHERE " + cond
		} else {
			query += " AND " + cond
		}
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY name LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Slug, &p.Description, &p.Unit, &p.Price, &p.StockQuantity,
			&p.ImageURL, &p.IsAvailable, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Update(ctx context.Context, id string, p *model.Product) (*model.Product, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE products
		SET category_id = $1, name = $2, slug = $3, description = $4, unit = $5, price = $6,
			stock_quantity = $7, image_url = $8, is_available = $9
		WHERE id = $10
		RETURNING `+productColumns,
		p.CategoryID, p.Name, p.Slug, p.Description, p.Unit, p.Price, p.StockQuantity, p.ImageURL, p.IsAvailable, id)

	updated, err := scanProduct(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: a product with this slug already exists", apperror.ErrConflict)
		}
		if apperror.IsForeignKeyViolation(err) {
			return nil, fmt.Errorf("%w: category does not exist", apperror.ErrValidation)
		}
		if apperror.IsCheckViolation(err) {
			return nil, fmt.Errorf("%w: %v", apperror.ErrValidation, err)
		}
		return nil, err
	}
	return updated, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		if apperror.IsForeignKeyViolation(err) {
			return fmt.Errorf("%w: product is still referenced by existing orders", apperror.ErrConflict)
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
