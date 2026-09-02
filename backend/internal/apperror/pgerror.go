package apperror

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Postgres SQLSTATE codes: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
)

func IsUniqueViolation(err error) bool {
	return pgErrorCode(err) == pgUniqueViolation
}

func IsForeignKeyViolation(err error) bool {
	return pgErrorCode(err) == pgForeignKeyViolation
}

func IsCheckViolation(err error) bool {
	return pgErrorCode(err) == pgCheckViolation
}

func pgErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
