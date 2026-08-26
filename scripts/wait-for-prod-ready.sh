#!/usr/bin/env bash
# Wait until production API reports ready (post-deploy smoke).
set -euo pipefail

URL="${1:-https://api.tikhvin-palomnik.ru/health/ready}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-36}"
SLEEP_SEC="${SLEEP_SEC:-10}"

echo "Waiting for ${URL} ..."
for attempt in $(seq 1 "$MAX_ATTEMPTS"); do
  code="$(curl -s -o /dev/null -w '%{http_code}' "$URL" || true)"
  if [[ "$code" == "200" ]]; then
    echo "Ready (HTTP 200) after ${attempt} attempt(s)."
    exit 0
  fi
  echo "Attempt ${attempt}/${MAX_ATTEMPTS}: HTTP ${code:-000}"
  sleep "$SLEEP_SEC"
done

echo "Production health check failed: ${URL}"
exit 1
