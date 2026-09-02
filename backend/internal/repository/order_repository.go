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

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

const orderColumns = `id, user_id, delivery_area_id, delivery_slot_id, delivery_date, delivery_address,
	status, total_amount::text, created_at, updated_at`

const orderItemColumns = "id, order_id, product_id, quantity::text, unit_price::text, subtotal::text, created_at"

// OrderItemInput is one requested line item; unit_price/subtotal are never trusted from the
// caller — they're computed inside the transaction from the live product row.
type OrderItemInput struct {
	ProductID string
	Quantity  string
}

func scanOrder(row pgx.Row) (*model.Order, error) {
	var o model.Order
	err := row.Scan(&o.ID, &o.UserID, &o.DeliveryAreaID, &o.DeliverySlotID, &o.DeliveryDate, &o.DeliveryAddress,
		&o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func scanOrderItem(row pgx.Row) (*model.OrderItem, error) {
	var oi model.OrderItem
	err := row.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity, &oi.UnitPrice, &oi.Subtotal, &oi.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &oi, nil
}

// Create places an order: locks + decrements stock for every item, snapshots the current
// product price as unit_price, and computes total_amount server-side. All-or-nothing.
func (r *OrderRepository) Create(ctx context.Context, o *model.Order, items []OrderItemInput) (*model.Order, []model.OrderItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	var orderID string
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, delivery_area_id, delivery_slot_id, delivery_date, delivery_address, total_amount)
		VALUES ($1, $2, $3, $4, $5, 0)
		RETURNING id
	`, o.UserID, o.DeliveryAreaID, o.DeliverySlotID, o.DeliveryDate, o.DeliveryAddress).Scan(&orderID)
	if err != nil {
		if apperror.IsForeignKeyViolation(err) {
			return nil, nil, fmt.Errorf("%w: user, delivery area, or delivery slot does not exist", apperror.ErrValidation)
		}
		return nil, nil, err
	}

	orderItems := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		var unitPrice string
		err := tx.QueryRow(ctx, `
			UPDATE products
			SET stock_quantity = stock_quantity - $1::numeric
			WHERE id = $2 AND is_available = true AND stock_quantity >= $1::numeric
			RETURNING price::text
		`, item.Quantity, item.ProductID).Scan(&unitPrice)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil, fmt.Errorf("%w: product %s is unavailable or does not have enough stock", apperror.ErrValidation, item.ProductID)
			}
			if apperror.IsCheckViolation(err) {
				return nil, nil, fmt.Errorf("%w: invalid quantity for product %s", apperror.ErrValidation, item.ProductID)
			}
			return nil, nil, err
		}

		row := tx.QueryRow(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3::numeric, $4::numeric, ($3::numeric * $4::numeric))
			RETURNING `+orderItemColumns, orderID, item.ProductID, item.Quantity, unitPrice)

		oi, err := scanOrderItem(row)
		if err != nil {
			if apperror.IsUniqueViolation(err) {
				return nil, nil, fmt.Errorf("%w: product %s appears more than once in the order", apperror.ErrValidation, item.ProductID)
			}
			if apperror.IsCheckViolation(err) {
				return nil, nil, fmt.Errorf("%w: invalid quantity for product %s", apperror.ErrValidation, item.ProductID)
			}
			return nil, nil, err
		}
		orderItems = append(orderItems, *oi)
	}

	row := tx.QueryRow(ctx, `
		UPDATE orders
		SET total_amount = (SELECT COALESCE(SUM(subtotal), 0) FROM order_items WHERE order_id = $1)
		WHERE id = $1
		RETURNING `+orderColumns, orderID)

	order, err := scanOrder(row)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return order, orderItems, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, []model.OrderItem, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+orderColumns+` FROM orders WHERE id = $1`, id)
	order, err := scanOrder(row)
	if err != nil {
		return nil, nil, err
	}

	items, err := r.itemsByOrderID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return order, items, nil
}

func (r *OrderRepository) itemsByOrderID(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+orderItemColumns+` FROM order_items WHERE order_id = $1 ORDER BY created_at`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.OrderItem{}
	for rows.Next() {
		var oi model.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity, &oi.UnitPrice, &oi.Subtotal, &oi.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, oi)
	}
	return items, rows.Err()
}

func (r *OrderRepository) List(ctx context.Context, userID *string, limit, offset int) ([]model.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders`
	var args []any
	if userID != nil {
		args = append(args, *userID)
		query += fmt.Sprintf(" WHERE user_id = $%d", len(args))
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []model.Order{}
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.DeliveryAreaID, &o.DeliverySlotID, &o.DeliveryDate, &o.DeliveryAddress,
			&o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// UpdateStatus performs a plain unconditional status update. The service layer is
// responsible for rejecting transitions out of a terminal (delivered/cancelled) state.
func (r *OrderRepository) UpdateStatus(ctx context.Context, id, status string) (*model.Order, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE orders SET status = $1 WHERE id = $2
		RETURNING `+orderColumns, status, id)

	order, err := scanOrder(row)
	if err != nil {
		if apperror.IsCheckViolation(err) {
			return nil, fmt.Errorf("%w: invalid status value", apperror.ErrValidation)
		}
		return nil, err
	}
	return order, nil
}
