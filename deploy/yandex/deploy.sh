#!/usr/bin/env bash
# Deploy this repo to production: https://tikhvin-palomnik.ru
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
REMOTE="${DEPLOY_REMOTE:-smailikin70@93.77.165.81}"
DEPLOY_DIR="${DEPLOY_DIR:-/opt/palomnik}"
SSH_KEY="${SSH_KEY:-}"
if [[ -z "$SSH_KEY" ]]; then
  shopt -s nullglob
  yc_keys=("${HOME}/.ssh/palomnik_yc" "${HOME}"/.ssh/*_yc)
  shopt -u nullglob
  for candidate in "${yc_keys[@]}" "${HOME}/.ssh/id_ed25519" "${HOME}/.ssh/id_rsa"; do
    if [[ -n "$candidate" && -f "$candidate" ]]; then
      SSH_KEY="$candidate"
      break
    fi
  done
fi
SSH_OPTS=(-o StrictHostKeyChecking=accept-new -o ConnectTimeout=15)
if [[ -n "$SSH_KEY" && -f "$SSH_KEY" ]]; then
  SSH_OPTS+=(-o IdentitiesOnly=yes -i "$SSH_KEY")
fi

COMPOSE="sudo docker compose --env-file .env.production -f docker-compose.yml -f docker-compose.prod.yml"

# Remote .env.production is the live secret file (Telegram, Cloudflare Worker, OAuth).
# Do not rsync it from the laptop/agent workspace — that would overwrite rotated tokens.
if [[ "${DEPLOY_CI:-}" != "1" ]] && [[ ! -f "$ROOT_DIR/.env.production" ]]; then
  echo "Missing local $ROOT_DIR/.env.production (needed as a reminder that prod secrets exist)."
  echo "Remote keeps its own .env.production; this file is not uploaded."
  exit 1
fi

echo "==> Production https://tikhvin-palomnik.ru"

echo "Waiting for SSH..."
ok=0
ssh_err=""
for i in $(seq 1 24); do
  if ssh_err="$(ssh "${SSH_OPTS[@]}" -o BatchMode=yes "$REMOTE" "echo ok" 2>&1)"; then
    ok=1
    break
  fi
  if [[ "$i" == "1" ]]; then
    echo "SSH first attempt failed ($REMOTE):"
    echo "$ssh_err" | tail -n 20
  fi
  sleep 5
done
if [[ "$ok" != "1" ]]; then
  echo "SSH failed: $REMOTE"
  echo "$ssh_err" | tail -n 20
  exit 1
fi

if ! ssh "${SSH_OPTS[@]}" "$REMOTE" "test -d '${DEPLOY_DIR}'"; then
  detected="$(ssh "${SSH_OPTS[@]}" "$REMOTE" "ls -d /opt/*/docker-compose.yml 2>/dev/null | sed 's|/docker-compose.yml||' | head -1" || true)"
  detected="${detected//$'\r'/}"
  detected="$(echo "$detected" | tr -d '\n')"
  if [[ -n "$detected" ]]; then
    echo "Using existing remote directory: $detected (set DEPLOY_DIR to override)"
    DEPLOY_DIR="$detected"
  else
    echo "Remote directory ${DEPLOY_DIR} is missing. Create it or set DEPLOY_DIR."
    exit 1
  fi
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
  --exclude .env.production \
  --exclude .env \
  --exclude frontend/.env.local \
  --exclude "data/uploads" \
  --exclude "backend/data/uploads" \
  "$ROOT_DIR/" "${REMOTE}:${DEPLOY_DIR}/"

# Never `compose down -v` — the named Postgres volume must survive.
# Disk is tight: rebuild frontend only; recreate api/worker from existing images + updated env.
echo "Rebuilding frontend + api/worker, then recreating stack (no postgres wipe)..."
ssh "${SSH_OPTS[@]}" "$REMOTE" "cd '${DEPLOY_DIR}' && \
  $COMPOSE build frontend api worker && \
  $COMPOSE up -d --no-deps frontend && \
  $COMPOSE up -d && \
  $COMPOSE ps"


echo "Done. https://tikhvin-palomnik.ru"
