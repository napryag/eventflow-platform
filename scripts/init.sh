#!/usr/bin/env sh

# init.sh currently prepares local environment files
# and provides a single entry point for future
# infrastructure initialization steps.

set -e

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
until docker compose ps postgres | grep -q healthy; do
  sleep 2
done

echo "Running Liquibase migrations..."
docker compose up -d liquibase
docker compose wait liquibase

echo ""
echo "Initialization completed."