#!/usr/bin/env bash
# Dumps the production database to a timestamped, gzipped file and prunes anything
# older than RETENTION_DAYS. Meant to be run on the server via cron — see the crontab
# line in docs/DEPLOYMENT.md. Restore a dump with scripts/restore-db.sh.
set -euo pipefail

cd "$(dirname "$0")/.."
COMPOSE="docker compose -f docker-compose.prod.yml"
BACKUP_DIR="${BACKUP_DIR:-./data/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-14}"

if [ -f .env ]; then
	set -a
	# shellcheck disable=SC1091
	source .env
	set +a
fi
: "${POSTGRES_USER:?Set POSTGRES_USER in .env first}"
: "${POSTGRES_DB:?Set POSTGRES_DB in .env first}"

mkdir -p "$BACKUP_DIR"
timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
out_file="$BACKUP_DIR/${POSTGRES_DB}_${timestamp}.sql.gz"

echo "==> Dumping $POSTGRES_DB to $out_file"
$COMPOSE exec -T db pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" | gzip > "$out_file"

echo "==> Pruning backups older than $RETENTION_DAYS days"
find "$BACKUP_DIR" -name "${POSTGRES_DB}_*.sql.gz" -mtime "+$RETENTION_DAYS" -delete

echo "Done: $out_file ($(du -h "$out_file" | cut -f1))"
