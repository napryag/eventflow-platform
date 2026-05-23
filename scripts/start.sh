#!/usr/bin/env sh

set -e

# This script provides a single entry point
# for managing local application services.

SCRIPT_TIMEOUT="${SCRIPT_TIMEOUT:-1800}" # 30 min

run_with_timeout() {
  if [ -z "$TIMEOUT_BIN" ]; then
    TIMEOUT_BIN=$(command -v timeout || command -v gtimeout)
  fi

  if [ -z "$TIMEOUT_BIN" ]; then
    echo "timeout utility not found"
    exit 1
  fi

  "$TIMEOUT_BIN" --foreground "$SCRIPT_TIMEOUT" "$@"
}

COMPOSE_FILE="docker-compose.yml"

SUPPORTED_SERVICES="authhub"

BUILD=false
BUILD_NO_CACHE=false
UP=false
DETACHED=false
STOP=false
RESTART=false
LOGS=false

TARGETS="all"

print_usage() {
  echo "Usage:"
  echo "  ./scripts/start.sh [options]"
  echo ""
  echo "Options:"
  echo "  -b              Build application service images"
  echo "  -B              Build application service images without cache"
  echo "  -u              Start application services"
  echo "  -d              Detached mode (only with -u)"
  echo "  -s              Stop application services"
  echo "  -r              Restart application services"
  echo "  -l              Show service logs"
  echo ""
  echo "Targets:"
  echo "  -t all|<list>   Service targets separated by comma (default: all)"
  echo ""
  echo "Supported targets:"
  echo "  authhub"
  echo "  all"
}

validate_service() {
  service="$1"

  case "$service" in
    authhub)
      ;;
    *)
      echo "Error: Unsupported service target '$service'"
      echo "Supported targets: authhub, all"
      exit 1
      ;;
  esac
}

resolve_targets() {

  if [ "$TARGETS" = "all" ]; then
    SERVICES="authhub"
    return
  fi

  SERVICES=$(echo "$TARGETS" | tr ',' ' ')

  for service in $SERVICES; do
    validate_service "$service"
  done
}

log_action() {
  echo ""
  echo "=================================================="
  echo "$1"
  echo "=================================================="
}

while getopts "bBudsrlt:" opt; do
  case "$opt" in
    b)
      BUILD=true
      ;;
    B)
      BUILD_NO_CACHE=true
      ;;
    u)
      UP=true
      ;;
    d)
      DETACHED=true
      ;;
    s)
      STOP=true
      ;;
    r)
      RESTART=true
      ;;
    l)
      LOGS=true
      ;;
    t)
      TARGETS="$OPTARG"
      ;;
    *)
      print_usage
      exit 1
      ;;
  esac
done

if [ "$DETACHED" = true ] && [ "$UP" != true ]; then
  echo "Error: -d can only be used together with -u"
  exit 1
fi

resolve_targets

if [ "$BUILD" = true ]; then
  log_action "Building application services"

  for service in $SERVICES; do
    echo "Building service: $service"
    run_with_timeout \
      docker compose -f "$COMPOSE_FILE" build "$service"
  done
fi

if [ "$BUILD_NO_CACHE" = true ]; then
  log_action "Building application services without Docker cache"

  for service in $SERVICES; do
    echo "Building service without cache: $service"
    run_with_timeout \
      docker compose -f "$COMPOSE_FILE" build --no-cache "$service"
  done
fi

if [ "$UP" = true ]; then
  log_action "Starting application services"

  if [ "$DETACHED" = true ]; then
    echo "Running services in detached mode"
    docker compose -f "$COMPOSE_FILE" up -d $SERVICES
  else
    echo "Running services in foreground mode"
    docker compose -f "$COMPOSE_FILE" up $SERVICES
  fi
fi

if [ "$STOP" = true ]; then
  log_action "Stopping application services"

  for service in $SERVICES; do
    echo "Stopping service: $service"
    run_with_timeout \
      docker compose -f "$COMPOSE_FILE" stop "$service"
  done
fi

if [ "$RESTART" = true ]; then
  log_action "Restarting application services"

  for service in $SERVICES; do
    echo "Restarting service: $service"
    run_with_timeout \
      docker compose -f "$COMPOSE_FILE" restart "$service"
  done
fi

if [ "$LOGS" = true ]; then
  log_action "Showing service logs"

  echo "Streaming logs for services: $SERVICES"

  docker compose -f "$COMPOSE_FILE" logs -f $SERVICES
fi