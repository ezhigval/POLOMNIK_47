#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
API_BASE="${API_BASE:-http://localhost:8080}"
CHECK_WORKER="${CHECK_WORKER:-1}"
MAX_FAILED="${MAX_FAILED:-0}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-40}"
SLEEP_SECONDS="${SLEEP_SECONDS:-3}"

if [[ -z "${ADMIN_TOKEN:-}" && -f "$ROOT_DIR/.env.production" ]]; then
  ADMIN_TOKEN="$(grep -E '^ADMIN_TOKEN=' "$ROOT_DIR/.env.production" | tail -1 | cut -d= -f2-)"
fi
ADMIN_TOKEN="${ADMIN_TOKEN:-dev-admin-token}"

docker_bin() {
  if docker info >/dev/null 2>&1; then
    docker "$@"
  else
    sudo docker "$@"
  fi
}

api_get() {
  local path="$1"
  local extra_header="${2:-}"
  if [[ "${API_VIA_DOCKER:-}" == "1" ]]; then
    local cid
    cid="$(docker_bin ps --format '{{.Names}}' | grep -E -- '-api-1$' | head -1)"
    if [[ -z "$cid" ]]; then
      echo "No running api container found" >&2
      exit 1
    fi
    if [[ -n "$extra_header" ]]; then
      docker_bin exec "$cid" wget -qO- --header="$extra_header" "http://127.0.0.1:8080${path}"
    else
      docker_bin exec "$cid" wget -qO- "http://127.0.0.1:8080${path}"
    fi
  else
    if [[ -n "$extra_header" ]]; then
      curl -sf -H "$extra_header" "${API_BASE}${path}"
    else
      curl -sf "${API_BASE}${path}"
    fi
  fi
}

echo "Checking API readiness..."
api_get "/health/ready" >/dev/null
echo "API ready"

echo "Checking management system-info..."
INFO="$(api_get "/api/v1/management/system-info" "X-Admin-Token: ${ADMIN_TOKEN}")"
python3 - "$INFO" "$MAX_FAILED" <<'PY'
import json, sys
payload = json.loads(sys.argv[1])
data = payload.get("data") or {}
outbox = data.get("outbox") or {}
failed = int(outbox.get("failed") or 0)
pending = int(outbox.get("pending") or 0)
max_failed = int(sys.argv[2])
print(f"outbox pending={pending} failed={failed}")
if failed > max_failed:
    error = outbox.get("latest_failed_error") or "unknown"
    print(f"FAILED: outbox failed count {failed} exceeds max {max_failed}: {error}", file=sys.stderr)
    sys.exit(1)
PY

if [[ "$CHECK_WORKER" == "1" ]]; then
  echo "Checking worker container health..."
  for attempt in $(seq 1 "$MAX_ATTEMPTS"); do
    HEALTH="$(docker_bin ps --filter "name=worker" --format '{{.Names}} {{.Status}}' 2>/dev/null | awk '/worker/ {print; found=1} END{if(!found) exit 2}')" || HEALTH=""
    if echo "$HEALTH" | grep -qi 'healthy'; then
      echo "Worker healthy"
      break
    fi
    if [[ "$attempt" -eq "$MAX_ATTEMPTS" ]]; then
      echo "Worker health is '${HEALTH:-missing}', expected healthy" >&2
      docker_bin ps --filter "name=worker" >&2 || true
      exit 1
    fi
    echo "Waiting for worker ($attempt/$MAX_ATTEMPTS, status=${HEALTH:-missing})..."
    sleep "$SLEEP_SECONDS"
  done
fi

echo "Ops check passed"
