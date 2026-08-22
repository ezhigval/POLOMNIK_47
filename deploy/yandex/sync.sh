#!/usr/bin/env bash
# Alias: same as deploy.sh (single production host).
set -euo pipefail
exec "$(dirname "$0")/deploy.sh"
