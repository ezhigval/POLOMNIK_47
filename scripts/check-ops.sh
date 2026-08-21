#!/usr/bin/env bash
set -euo pipefail

API_BASE="${API_BASE:-http://localhost:8080}"
ADMIN_TOKEN="${ADMIN_TOKEN:-dev-admin-token}"
CHECK_WORKER="${CHECK_WORKER:-1}"
MAX_FAILED="${MAX_FAILED:-0}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-40}"
SLEEP_SECONDS="${SLEEP_SECONDS:-3}"

echo "Checking API readiness..."
curl -sf "${API_BASE}/health/ready" >/dev/null
echo "API ready"

echo "Checking management system-info..."
INFO="$(curl -sf -H "X-Admin-Token: ${ADMIN_TOKEN}" "${API_BASE}/api/v1/management/system-info")"
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
    HEALTH="$(docker compose ps --format '{{.Service}} {{.Health}}' 2>/dev/null | awk '$1=="worker" {print $2; found=1} END{if(!found) exit 2}')" || HEALTH=""
    if [[ "$HEALTH" == "healthy" ]]; then
      echo "Worker healthy"
      break
    fi
    if [[ "$attempt" -eq "$MAX_ATTEMPTS" ]]; then
      echo "Worker health is '${HEALTH:-missing}', expected healthy" >&2
      docker compose ps worker >&2 || true
      exit 1
    fi
    echo "Waiting for worker ($attempt/$MAX_ATTEMPTS, health=${HEALTH:-missing})..."
    sleep "$SLEEP_SECONDS"
  done
fi

echo "Ops check passed"
