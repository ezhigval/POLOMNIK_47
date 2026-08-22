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

POSTGRES_USER="${POSTGRES_USER:-polomnik}"
POSTGRES_DB="${POSTGRES_DB:-polomnik}"

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
OUTPUT="$BACKUP_DIR/polomnik-${TIMESTAMP}.sql.gz"

echo "Creating backup from $CONTAINER: $OUTPUT"
docker_bin exec "$CONTAINER" pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --no-owner --no-acl | gzip > "$OUTPUT"

find "$BACKUP_DIR" -name 'polomnik-*.sql.gz' -mtime +"$RETENTION_DAYS" -delete 2>/dev/null || true

echo "Backup complete ($(du -h "$OUTPUT" | cut -f1))"
