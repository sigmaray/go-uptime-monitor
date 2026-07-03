#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

if [[ -z "${GO_UPTIME_MONITOR_TEST_DATABASE_NAME:-}" ]]; then
  echo "GO_UPTIME_MONITOR_TEST_DATABASE_NAME is not set" >&2
  exit 1
fi

export GO_UPTIME_MONITOR_ENABLE_PLAYWRIGHT_API=true
export GO_UPTIME_MONITOR_DATABASE_NAME="$GO_UPTIME_MONITOR_TEST_DATABASE_NAME"

compose() {
  docker compose -p e2e-uptime -f docker-compose.with-infra.yml "$@"
}

compose build app
compose up -d postgres --wait
compose run --rm app ./uptime-monitor migrate
compose up --build
