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

  if [ -f .env.example ]; then
    cp .env.example .env
    echo ".env created from .env.example"
  else
    touch .env
    echo "Empty .env file created"
  fi
else
  echo ".env already exists"
fi

echo ""
echo "Initialization completed."
echo "Initialized:"
echo " - Docker availability check"
echo " - Local environment file (.env)"