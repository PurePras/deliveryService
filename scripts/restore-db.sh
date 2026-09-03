#!/usr/bin/env bash
# Restores a .sql.gz dump produced by scripts/backup-db.sh. DESTRUCTIVE: drops and
# recreates every object in the target database first.
#
# Usage: scripts/restore-db.sh ./data/backups/shri_ram_service_20260101T000000Z.sql.gz
set -euo pipefail

cd "$(dirname "$0")/.."
COMPOSE="docker compose -f docker-compose.prod.yml"

dump_file="${1:?Usage: $0 <path-to-dump.sql.gz>}"
[ -f "$dump_file" ] || { echo "No such file: $dump_file" >&2; exit 1; }

if [ -f .env ]; then
	set -a
	# shellcheck disable=SC1091
	source .env
	set +a
fi
: "${POSTGRES_USER:?Set POSTGRES_USER in .env first}"
: "${POSTGRES_DB:?Set POSTGRES_DB in .env first}"

read -rp "This will WIPE the current '$POSTGRES_DB' database and replace it with $dump_file. Type 'yes' to continue: " confirm
[ "$confirm" = "yes" ] || { echo "Aborted."; exit 1; }

echo "==> Stopping the backend so nothing writes during restore"
$COMPOSE stop backend

echo "==> Dropping and recreating $POSTGRES_DB"
$COMPOSE exec -T db psql -U "$POSTGRES_USER" -d postgres -c "DROP DATABASE IF EXISTS \"$POSTGRES_DB\";"
$COMPOSE exec -T db psql -U "$POSTGRES_USER" -d postgres -c "CREATE DATABASE \"$POSTGRES_DB\";"

echo "==> Restoring $dump_file"
gunzip -c "$dump_file" | $COMPOSE exec -T db psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"

echo "==> Restarting the backend"
$COMPOSE start backend

echo "Done."
