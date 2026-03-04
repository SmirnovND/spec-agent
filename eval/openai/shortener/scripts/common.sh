#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

ARTIFACT_DIR="${ARTIFACT_DIR:-$WORKSPACE/.eval/openai/artifacts}"
mkdir -p "$ARTIFACT_DIR"

SERVER_CONTAINER_NAME="eval-shortener-server"
SERVER_IMAGE_NAME="eval-shortener-server:latest"

prepare_server_config() {
  mkdir -p "$WORKSPACE/cmd/server"
  cat > "$WORKSPACE/cmd/server/config.yaml" <<'YAML'
db:
  dsn: "postgresql://developer:developer@localhost:5432/gobase?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5

app:
  run_addr: "0.0.0.0:8080"
  shortener_host: "http://localhost:8080"

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
YAML
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
  (
    cd "$WORKSPACE"
    make migrate-up
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

wait_http_ready() {
  local retries=40
  local i
  for i in $(seq 1 "$retries"); do
    if curl -fsS "http://localhost:8080/ping" >/dev/null 2>&1; then
      return
    fi
    sleep 1
  done

  docker logs "$SERVER_CONTAINER_NAME" > "$ARTIFACT_DIR/server.log" 2>&1 || true
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
