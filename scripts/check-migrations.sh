#!/usr/bin/env bash
# CI guard: INSERT/UPDATE migrations with semicolons inside SQL must use goose StatementBegin.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS="$ROOT/backend/migrations"
failed=0

for file in "$MIGRATIONS"/*.sql; do
  [[ -f "$file" ]] || continue
  base="$(basename "$file")"

  if rg -q 'INSERT INTO|UPDATE ' "$file" && rg -q ';' "$file"; then
    if ! rg -q 'StatementBegin' "$file"; then
      echo "FAIL $base: SQL with semicolons must use -- +goose StatementBegin/StatementEnd (see 00025)"
      failed=1
    fi
  fi

  if [[ "$base" =~ ^0+([0-9]+)_ ]]; then
    num="${BASH_REMATCH[1]}"
    num=$((10#$num))
    if (( num < 1 || num > 999 )); then
      echo "FAIL $base: unexpected migration number"
      failed=1
    fi
  else
    echo "FAIL $base: filename must match NNNNN_name.sql"
    failed=1
  fi
done

# Ensure sequential numbering without gaps (00001 … 000NN).
max=0
for file in "$MIGRATIONS"/*.sql; do
  [[ -f "$file" ]] || continue
  base="$(basename "$file")"
  num="${base%%_*}"
  num=$((10#$num))
  if (( num > max )); then max=$num; fi
done

for (( i=1; i<=max; i++ )); do
  printf -v want "%05d" "$i"
  if ! compgen -G "$MIGRATIONS/${want}_*.sql" > /dev/null; then
    echo "FAIL: missing migration ${want}_*.sql (gap before max $max)"
    failed=1
    break
  fi
done

if (( failed != 0 )); then
  exit 1
fi

echo "Migrations OK ($max files, StatementBegin rule checked)"
