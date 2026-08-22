#!/usr/bin/env bash
# Deploy this repo to production: https://tikhvin-palomnik.ru
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
REMOTE="${DEPLOY_REMOTE:-smailikin70@93.77.165.81}"
SSH_KEY="${SSH_KEY:-${HOME}/.ssh/polomnik_yc}"
SSH_OPTS=(-o StrictHostKeyChecking=accept-new -o ConnectTimeout=15 -o IdentitiesOnly=yes)
if [[ -f "$SSH_KEY" ]]; then
  SSH_OPTS+=(-i "$SSH_KEY")
fi

COMPOSE="sudo docker compose --env-file .env.production -f docker-compose.yml -f docker-compose.prod.yml"

if [[ ! -f "$ROOT_DIR/.env.production" ]]; then
  echo "Missing $ROOT_DIR/.env.production"
  exit 1
fi

echo "==> Production https://tikhvin-palomnik.ru"

echo "Waiting for SSH..."
ok=0
for _ in $(seq 1 24); do
  if ssh "${SSH_OPTS[@]}" "$REMOTE" "echo ok" >/dev/null 2>&1; then
    ok=1
    break
  fi
  sleep 5
done
if [[ "$ok" != "1" ]]; then
  echo "SSH failed: $REMOTE"
  exit 1
fi

echo "Syncing..."
rsync -az --delete \
  -e "ssh ${SSH_OPTS[*]}" \
  --exclude .git \
  --exclude .idea \
  --exclude .DS_Store \
  --exclude node_modules \
  --exclude .next \
  --exclude backups \
  --exclude frontend/.env.local \
  --exclude "data/uploads" \
  --exclude "backend/data/uploads" \
  "$ROOT_DIR/" "${REMOTE}:/opt/polomnik/"

# Frontend + Caddy change often; a full stack --build also rebuilds goose/api and fills the 20 GB disk.
echo "Rebuilding frontend, then applying prod overlay (unpublished DB/API ports)..."
ssh "${SSH_OPTS[@]}" "$REMOTE" "cd /opt/polomnik && \
  $COMPOSE up -d --no-deps --build frontend && \
  $COMPOSE up -d --no-build && \
  $COMPOSE ps"

echo "Done. https://tikhvin-palomnik.ru"
