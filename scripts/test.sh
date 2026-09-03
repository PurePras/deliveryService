#!/usr/bin/env bash
# Runs the backend Go test suite. Unit tests always run; the integration suite
# (against the real dockerized Postgres) only runs when DATABASE_URL is set.
set -euo pipefail

cd "$(dirname "$0")/../backend"

echo "==> Unit tests (go test ./...)"
go test ./...

if [ -z "${DATABASE_URL:-}" ]; then
	echo "==> DATABASE_URL is not set; skipping integration tests"
	echo "    (export it — see backend/.env.example — to also run the tests in backend/tests/)"
	exit 0
fi

echo "==> Integration tests (go test -tags=integration ./tests/...)"
go test -tags=integration ./tests/...
