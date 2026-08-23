#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-$ROOT_DIR/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"

if [[ -f "$ROOT_DIR/.env.production" ]]; then
  POSTGRES_USER="$(grep -E '^POSTGRES_USER=' "$ROOT_DIR/.env.production" | tail -1 | cut -d= -f2-)"
  POSTGRES_DB="$(grep -E '^POSTGRES_DB=' "$ROOT_DIR/.env.production" | tail -1 | cut -d= -f2-)"
fi

POSTGRES_USER="${POSTGRES_USER:-palomnik}"
POSTGRES_DB="${POSTGRES_DB:-palomnik}"

docker_bin() {
  if docker info >/dev/null 2>&1; then
    docker "$@"
  else
    sudo docker "$@"
  fi
}

resolve_container() {
  if [[ -n "${POSTGRES_CONTAINER:-}" ]]; then
    echo "$POSTGRES_CONTAINER"
    return
  fi
  docker_bin ps --format '{{.Names}}' | grep -E 'postgres-1$' | head -1
}

CONTAINER="$(resolve_container)"
if [[ -z "$CONTAINER" ]]; then
  echo "No running postgres container found" >&2
  exit 1
fi

mkdir -p "$BACKUP_DIR"
OUTPUT="$BACKUP_DIR/palomnik-${TIMESTAMP}.sql.gz"

echo "Creating backup from $CONTAINER: $OUTPUT"
docker_bin exec "$CONTAINER" pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --no-owner --no-acl | gzip > "$OUTPUT"

find "$BACKUP_DIR" -name '*.sql.gz' -mtime +"$RETENTION_DAYS" -delete 2>/dev/null || true

API_CONTAINER="${API_CONTAINER:-}"
if [[ -z "$API_CONTAINER" ]]; then
  API_CONTAINER="$(docker_bin ps --format '{{.Names}}' | grep -E 'api-1$' | head -1 || true)"
fi

if [[ -n "$API_CONTAINER" ]]; then
  docker_bin cp "$OUTPUT" "$API_CONTAINER:/tmp/$(basename "$OUTPUT")"
  if docker_bin exec "$API_CONTAINER" /app/polomnik-backup-offsite "/tmp/$(basename "$OUTPUT")"; then
    echo "Backup status recorded (offsite if S3 keys are set)"
  else
    echo "Offsite/status step failed (local dump is kept)" >&2
  fi
  docker_bin exec "$API_CONTAINER" rm -f "/tmp/$(basename "$OUTPUT")" || true
fi

echo "Backup complete ($(du -h "$OUTPUT" | cut -f1))"
