package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
)

// OTPCode is a row in otp_codes. It never crosses into a dto/API response — the
// service layer is the only consumer, deciding expiry/attempt/consumed checks itself
// so that logic lives in one place instead of being duplicated in SQL.
type OTPCode struct {
	ID         string
	Phone      string
	CodeHash   string
	ExpiresAt  time.Time
	Attempts   int
	ConsumedAt *time.Time
	CreatedAt  time.Time
}

type OTPRepository struct {
	pool *pgxpool.Pool
}

func NewOTPRepository(pool *pgxpool.Pool) *OTPRepository {
	return &OTPRepository{pool: pool}
}

const otpColumns = "id, phone, code_hash, expires_at, attempts, consumed_at, created_at"

func scanOTPCode(row pgx.Row) (*OTPCode, error) {
	var o OTPCode
	err := row.Scan(&o.ID, &o.Phone, &o.CodeHash, &o.ExpiresAt, &o.Attempts, &o.ConsumedAt, &o.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *OTPRepository) Create(ctx context.Context, phone, codeHash string, expiresAt time.Time) (*OTPCode, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO otp_codes (phone, code_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING `+otpColumns, phone, codeHash, expiresAt)
	return scanOTPCode(row)
}

// GetLatestByPhone returns the most recently created OTP for phone, regardless of
// whether it's since been consumed or expired — callers decide what that means.
// Only ever one "live" code per phone: requesting a new one supersedes the last.
func (r *OTPRepository) GetLatestByPhone(ctx context.Context, phone string) (*OTPCode, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT `+otpColumns+`
		FROM otp_codes
		WHERE phone = $1
		ORDER BY created_at DESC
		LIMIT 1`, phone)
	return scanOTPCode(row)
}

// CountSince reports how many OTPs have been generated for phone since the given
// time, for the hourly-cap check in AuthService.RequestLoginOTP.
func (r *OTPRepository) CountSince(ctx context.Context, phone string, since time.Time) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM otp_codes WHERE phone = $1 AND created_at > $2`, phone, since,
	).Scan(&count)
	return count, err
}

func (r *OTPRepository) IncrementAttempts(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE otp_codes SET attempts = attempts + 1 WHERE id = $1`, id)
	return err
}

func (r *OTPRepository) MarkConsumed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE otp_codes SET consumed_at = now() WHERE id = $1`, id)
	return err
}
