#!/usr/bin/env bash
set -euo pipefail

URL="${FRONTEND_URL:-http://localhost:3000}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-60}"
SLEEP_SECONDS="${SLEEP_SECONDS:-2}"

for attempt in $(seq 1 "$MAX_ATTEMPTS"); do
  if curl -sf "$URL" >/dev/null; then
    echo "Frontend ready at $URL"
    exit 0
  fi
  echo "Waiting for frontend ($attempt/$MAX_ATTEMPTS)..."
  sleep "$SLEEP_SECONDS"
done

echo "Frontend did not become ready at $URL"
exit 1
