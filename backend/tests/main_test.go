//go:build integration

// Package tests holds integration tests that run against a real Postgres database
// (the one started by docker-compose.yml) rather than mocks — the repository layer
// isn't behind interfaces, and refactoring it just for mockability isn't justified
// here. Run via `go test -tags=integration ./tests/...`, or scripts/test.sh which
// skips this package cleanly when DATABASE_URL isn't set.
package tests

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PurePras/shri-ram-service/backend/internal/database"
)

// testPool is shared read-write access to the real database for every test in this
// package, set up once in TestMain.
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL is not set; skipping integration tests (see scripts/test.sh)")
		os.Exit(0)
	}

	pool, err := database.NewPool(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("connect to database:", err)
		os.Exit(1)
	}
	testPool = pool

	code := m.Run()
	pool.Close()
	os.Exit(code)
}

// randomSuffix gives test-created rows (categories, products, delivery areas...) a
// name/slug/pincode that won't collide with seed data or a previous run's leftovers.
func randomSuffix() string {
	return fmt.Sprintf("%d", rand.N(1_000_000))
}

// randomTestPhone returns a syntactically valid (see service.requirePhone), unique
// 10-digit test phone number.
func randomTestPhone() string {
	return fmt.Sprintf("9%09d", rand.N(1_000_000_000))
}
