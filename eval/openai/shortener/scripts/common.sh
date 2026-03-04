#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

ARTIFACT_DIR="${ARTIFACT_DIR:-$WORKSPACE/.eval/openai/artifacts}"
mkdir -p "$ARTIFACT_DIR"

SERVER_CONTAINER_NAME="eval-shortener-server"
SERVER_IMAGE_NAME="eval-shortener-server:latest"
SERVER_MODE="${SERVER_MODE:-docker}" # docker | process
SERVER_PID_FILE="$ARTIFACT_DIR/server.pid"
DB_HOST="${DB_HOST:-localhost}"
RABBITMQ_HOST="${RABBITMQ_HOST:-localhost}"
APP_HOST="${APP_HOST:-http://localhost:8080}"

prepare_server_config() {
  mkdir -p "$WORKSPACE/cmd/server"
  cat > "$WORKSPACE/cmd/server/config.yaml" <<'YAML'
db:
  dsn: "__DB_DSN__"
  max_open_conns: 25
  max_idle_conns: 5

app:
  run_addr: "0.0.0.0:8080"
  shortener_host: "__APP_HOST__"

rabbitmq:
  url: "amqp://guest:guest@__RABBITMQ_HOST__:5672/"
YAML
  local db_dsn="postgresql://developer:developer@${DB_HOST}:5432/gobase?sslmode=disable"
  sed -i "s|__DB_DSN__|$db_dsn|g" "$WORKSPACE/cmd/server/config.yaml"
  sed -i "s|__APP_HOST__|$APP_HOST|g" "$WORKSPACE/cmd/server/config.yaml"
  sed -i "s|__RABBITMQ_HOST__|$RABBITMQ_HOST|g" "$WORKSPACE/cmd/server/config.yaml"
}

install_migrate_if_needed() {
  if command -v migrate >/dev/null 2>&1; then
    return
  fi
  GOBIN="${GOBIN:-$HOME/go/bin}" go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0
  export PATH="$GOBIN:$PATH"
}

apply_migrations() {
  install_migrate_if_needed
  local db_dsn="postgresql://developer:developer@${DB_HOST}:5432/gobase?sslmode=disable"
  (
    cd "$WORKSPACE"
    make migrate-up DB_DSN="$db_dsn"
  )
}

build_server_image() {
  cat > "$WORKSPACE/Dockerfile.eval" <<'DOCKERFILE'
FROM golang:1.24 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /out/server /app/server
COPY cmd/server/config.yaml /app/config.yaml
EXPOSE 8080
ENTRYPOINT ["/app/server", "/app/config.yaml"]
DOCKERFILE

  docker build -f "$WORKSPACE/Dockerfile.eval" -t "$SERVER_IMAGE_NAME" "$WORKSPACE"
}

start_server_container() {
  stop_server_container || true
  docker run -d \
    --name "$SERVER_CONTAINER_NAME" \
    --network host \
    --add-host host.docker.internal:host-gateway \
    "$SERVER_IMAGE_NAME"
}

stop_server_container() {
  if docker ps -a --format '{{.Names}}' | grep -q "^${SERVER_CONTAINER_NAME}$"; then
    docker rm -f "$SERVER_CONTAINER_NAME" >/dev/null 2>&1 || true
  fi
}

start_server_process() {
  stop_server_process || true
  (
    cd "$WORKSPACE"
    nohup go run ./cmd/server/main.go ./cmd/server/config.yaml >"$ARTIFACT_DIR/server.log" 2>&1 &
    echo $! >"$SERVER_PID_FILE"
  )
}

stop_server_process() {
  if [ -f "$SERVER_PID_FILE" ]; then
    pid="$(cat "$SERVER_PID_FILE" || true)"
    if [ -n "${pid:-}" ] && kill -0 "$pid" >/dev/null 2>&1; then
      kill "$pid" >/dev/null 2>&1 || true
      sleep 1
      kill -9 "$pid" >/dev/null 2>&1 || true
    fi
    rm -f "$SERVER_PID_FILE"
  fi
}

start_server() {
  if [ "$SERVER_MODE" = "process" ]; then
    start_server_process
    return
  fi
  build_server_image
  start_server_container
}

stop_server() {
  stop_server_process || true
  stop_server_container || true
}

wait_http_ready() {
  local retries=40
  local i
  for i in $(seq 1 "$retries"); do
    if curl -fsS "http://localhost:8080/ping" >/dev/null 2>&1; then
      return
    fi
    sleep 1
  done

  if [ "$SERVER_MODE" = "docker" ]; then
    docker logs "$SERVER_CONTAINER_NAME" > "$ARTIFACT_DIR/server.log" 2>&1 || true
  fi
  echo "server did not become ready, logs saved to $ARTIFACT_DIR/server.log"
  return 1
}

extract_json_field() {
  local field="$1"
  local file="$2"
  jq -r "$field" "$file"
}

assert_equals() {
  local expected="$1"
  local actual="$2"
  local message="$3"
  if [ "$expected" != "$actual" ]; then
    echo "assert failed: $message expected='$expected' actual='$actual'"
    return 1
  fi
}

assert_prefix() {
  local prefix="$1"
  local value="$2"
  local message="$3"
  case "$value" in
    "$prefix"*) ;;
    *)
      echo "assert failed: $message expected_prefix='$prefix' actual='$value'"
      return 1
      ;;
  esac
}
