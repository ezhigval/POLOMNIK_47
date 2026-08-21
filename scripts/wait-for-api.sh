#!/usr/bin/env bash
set -euo pipefail

URL="${API_URL:-http://localhost:8080/api/v1/health/ready}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-60}"
SLEEP_SECONDS="${SLEEP_SECONDS:-2}"

for attempt in $(seq 1 "$MAX_ATTEMPTS"); do
  if curl -sf "$URL" >/dev/null; then
    echo "API ready at $URL"
    exit 0
  fi
  echo "Waiting for API ($attempt/$MAX_ATTEMPTS)..."
  sleep "$SLEEP_SECONDS"
done

echo "API did not become ready at $URL"
exit 1
