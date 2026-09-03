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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const userColumns = "id, name, phone, email, password_hash, role, created_at, updated_at"

func scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Phone, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// Create ignores any role on u — new users are always created as model.RoleCustomer.
// Promoting a user to admin is an intentionally out-of-band operation (direct DB access),
// never something a client request can trigger.
func (r *UserRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO users (name, phone, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+userColumns, u.Name, u.Phone, u.Email, u.PasswordHash, model.RoleCustomer)

	created, err := scanUser(row)
	if err != nil {
		if apperror.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%w: a user with this phone or email already exists", apperror.ErrConflict)
		}
		return nil, err
	}
	return created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE phone = $1`, phone)
	return scanUser(row)
}

func (r *UserRepository) UpdatePasswordHash(ctx context.Context, userID, passwordHash string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE users SET password_hash = $1 WHERE id = $2`, passwordHash, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
