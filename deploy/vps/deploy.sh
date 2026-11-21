#!/usr/bin/env bash
# Deploy the checked-out revision to an already provisioned VPS.
# Production secrets, uploads, and backups stay on the target server.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
VPS_SSH_HOST="${VPS_SSH_HOST:-}"
VPS_DEPLOY_DIR="${VPS_DEPLOY_DIR:-/opt/polomnik}"
SSH_KEY="${SSH_KEY:-}"
VPS_ENV_FILE="${VPS_ENV_FILE:-.env.production}"
if [[ ! "$VPS_ENV_FILE" =~ ^\.env\.production(\.[a-zA-Z0-9_-]+)?$ ]]; then
  echo "Invalid production environment filename." >&2
  exit 2
fi

if [[ -z "$VPS_SSH_HOST" ]]; then
  echo "Set VPS_SSH_HOST to a verified user@host value." >&2
  exit 2
fi
if [[ ! "$VPS_SSH_HOST" =~ ^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+$ ]]; then
  echo "VPS_SSH_HOST must be user@host (IPv6 and custom ports are not supported)." >&2
  exit 2
fi
if [[ "$VPS_DEPLOY_DIR" != /* || "$VPS_DEPLOY_DIR" == "/" ]]; then
  echo "VPS_DEPLOY_DIR must be an absolute, non-root path." >&2
  exit 2
fi
if [[ -z "$SSH_KEY" || ! -f "$SSH_KEY" ]]; then
  echo "Set SSH_KEY to the private key for the VPS deploy user." >&2
  exit 2
fi
if [[ ! -f "$ROOT_DIR/docker-compose.yml" || ! -f "$ROOT_DIR/docker-compose.prod.yml" ]]; then
  echo "Run this script from a POLOMNIK_47 checkout." >&2
  exit 2
fi

VPS_HOST="${VPS_SSH_HOST#*@}"
KNOWN_HOSTS="${KNOWN_HOSTS:-${HOME}/.ssh/known_hosts}"
if [[ ! -f "$KNOWN_HOSTS" ]] || ! ssh-keygen -F "$VPS_HOST" -f "$KNOWN_HOSTS" >/dev/null; then
  echo "The VPS host key is not pinned in $KNOWN_HOSTS." >&2
  echo "Verify the fingerprint with the provider, then add its host key before deploying." >&2
  exit 2
fi

SSH_OPTS=(-o BatchMode=yes -o IdentitiesOnly=yes -o StrictHostKeyChecking=yes \
  -o UserKnownHostsFile="$KNOWN_HOSTS" -o GlobalKnownHostsFile=/dev/null \
  -o ConnectTimeout=15 -o ServerAliveInterval=15 -o ServerAliveCountMax=3 -i "$SSH_KEY")
printf -v RSYNC_SSH '%q ' ssh "${SSH_OPTS[@]}"
printf -v REMOTE_DIR_Q '%q' "$VPS_DEPLOY_DIR"

COMPOSE="docker compose --env-file ${VPS_ENV_FILE} -f docker-compose.yml -f docker-compose.prod.yml -f deploy/vps/compose.prod.yml"

echo "==> Checking the existing VPS deployment directory"
ssh "${SSH_OPTS[@]}" "$VPS_SSH_HOST" \
  "test -d ${REMOTE_DIR_Q} && test -s ${REMOTE_DIR_Q}/${VPS_ENV_FILE} && docker compose version >/dev/null"

# Only this directory contains application source exclusively. Remove retired modules
# here so Next does not compile obsolete routes from the previous release.
rsync -az --delete-after -e "$RSYNC_SSH" "$ROOT_DIR/frontend/src/" "$VPS_SSH_HOST:$VPS_DEPLOY_DIR/frontend/src/"

echo "==> Syncing application files (secrets and persistent data are excluded)"
rsync -az -e "$RSYNC_SSH" \
  --exclude .git \
  --exclude .idea \
  --exclude .DS_Store \
  --exclude node_modules \
  --exclude .next \
  --exclude backups \
  --exclude '.env*' \
  --exclude .env.production.live \
  --exclude .env.mail \
  --exclude .env \
  --exclude frontend/.env.local \
  --exclude data/uploads \
  --exclude backend/data/uploads \
  "$ROOT_DIR/" "${VPS_SSH_HOST}:${VPS_DEPLOY_DIR}/"

revision="$(git -C "$ROOT_DIR" rev-parse HEAD)"
if [[ ! "$revision" =~ ^[0-9a-f]{40}$ ]]; then exit 2; fi
printf '%s\n' "$revision" | ssh "${SSH_OPTS[@]}" "$VPS_SSH_HOST" "cat > ${REMOTE_DIR_Q}/frontend/public/release.txt"

echo "==> Building and starting the Compose production stack"
ssh "${SSH_OPTS[@]}" "$VPS_SSH_HOST" \
  "cd ${REMOTE_DIR_Q} && ${COMPOSE} build frontend api worker mail migrate && ${COMPOSE} up -d && ${COMPOSE} ps"

for attempt in {1..30}; do
  deployed="$(curl -fsS --max-time 10 --resolve "tikhvin-palomnik.ru:443:$VPS_HOST" \
    "https://tikhvin-palomnik.ru/release.txt?revision=$revision" || true)"
  if [[ "$deployed" == "$revision" ]]; then
    echo "Selectel serves $revision."
    exit 0
  fi
  sleep 5
done
echo "The site did not serve the expected revision $revision." >&2
exit 1
