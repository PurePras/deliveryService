package apperror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestPgErrorPredicates(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantUnique bool
		wantFK     bool
		wantCheck  bool
	}{
		{
			name:       "unique violation",
			err:        &pgconn.PgError{Code: pgUniqueViolation},
			wantUnique: true,
		},
		{
			name:   "foreign key violation",
			err:    &pgconn.PgError{Code: pgForeignKeyViolation},
			wantFK: true,
		},
		{
			name:      "check violation",
			err:       &pgconn.PgError{Code: pgCheckViolation},
			wantCheck: true,
		},
		{
			name:       "wrapped unique violation is still detected through errors.As",
			err:        fmt.Errorf("insert user: %w", &pgconn.PgError{Code: pgUniqueViolation}),
			wantUnique: true,
		},
		{
			name: "unrelated pg error code matches none of the three",
			err:  &pgconn.PgError{Code: "42601"}, // syntax_error
		},
		{
			name: "plain non-pg error",
			err:  errors.New("boom"),
		},
		{
			name: "nil error",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUniqueViolation(tt.err); got != tt.wantUnique {
				t.Errorf("IsUniqueViolation = %v, want %v", got, tt.wantUnique)
			}
			if got := IsForeignKeyViolation(tt.err); got != tt.wantFK {
				t.Errorf("IsForeignKeyViolation = %v, want %v", got, tt.wantFK)
			}
			if got := IsCheckViolation(tt.err); got != tt.wantCheck {
				t.Errorf("IsCheckViolation = %v, want %v", got, tt.wantCheck)
			}
		})
	}
}
