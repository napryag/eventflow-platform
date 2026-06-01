#!/usr/bin/env sh

# init.sh currently prepares local environment files
# and provides a single entry point for future
# infrastructure initialization steps.

set -e

SCRIPT_TIMEOUT="${SCRIPT_TIMEOUT:-600}" # 10 min

run_with_timeout() {
  TIMEOUT_BIN=$(command -v timeout || command -v gtimeout)

  if [ -z "$TIMEOUT_BIN" ]; then
    echo "timeout utility not found"
    exit 1
  fi

  "$TIMEOUT_BIN" --foreground "$SCRIPT_TIMEOUT" "$@"
}

echo "Checking Docker availability..."

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is not installed or not available in PATH."
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "Docker daemon is not running."
  exit 1
fi

echo "Docker is available."

if [ ! -f .env ]; then
  echo ".env file not found."
  exit 1
fi

echo "Starting PostgreSQL..."
docker compose up -d postgres

echo "Waiting for PostgreSQL..."
run_with_timeout sh -c '
until docker compose ps postgres | grep -q healthy; do
  echo "PostgreSQL is not healthy yet..."
  sleep 2
done
'

echo "Running Liquibase migrations..."
docker compose up -d liquibase
run_with_timeout docker compose wait liquibase

echo ""
echo "Initialization completed."