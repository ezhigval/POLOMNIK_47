#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-$ROOT_DIR/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"

POSTGRES_USER="${POSTGRES_USER:-polomnik}"
POSTGRES_DB="${POSTGRES_DB:-polomnik}"
CONTAINER="${POSTGRES_CONTAINER:-polomnik_47-postgres-1}"

mkdir -p "$BACKUP_DIR"
OUTPUT="$BACKUP_DIR/polomnik-${TIMESTAMP}.sql.gz"

echo "Creating backup: $OUTPUT"
docker exec "$CONTAINER" pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --no-owner --no-acl | gzip > "$OUTPUT"

find "$BACKUP_DIR" -name 'polomnik-*.sql.gz' -mtime +"$RETENTION_DAYS" -delete 2>/dev/null || true

echo "Backup complete ($(du -h "$OUTPUT" | cut -f1))"
