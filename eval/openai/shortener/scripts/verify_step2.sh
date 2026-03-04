#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/common.sh"

prepare_server_config
apply_migrations
start_server
trap 'stop_server' EXIT
wait_http_ready

mkdir -p "$WORKSPACE/.eval/openai"

CREATE_PAYLOAD="$WORKSPACE/.eval/openai/step2_create_payload.json"
CREATE_RESPONSE="$WORKSPACE/.eval/openai/step2_create_response.json"

expired_at="$(date -u -d '1 minute ago' +%Y-%m-%dT%H:%M:%SZ)"

jq -n \
  --arg url "https://example.com/expired" \
  --arg expires_at "$expired_at" \
  '{url:$url, expires_at:$expires_at}' > "$CREATE_PAYLOAD"

curl -fsS \
  -X POST "http://localhost:8080/api/v1/short-links" \
  -H "Content-Type: application/json" \
  --data-binary "@$CREATE_PAYLOAD" \
  > "$CREATE_RESPONSE"

short_url="$(extract_json_field '.short_url' "$CREATE_RESPONSE")"
assert_prefix "${APP_HOST}/" "$short_url" "short_url host mismatch"

(
  cd "$WORKSPACE"
  go run ./cmd/crons/cleanup_short_links/main.go ./cmd/server/config.yaml
)

status_code="$(curl -sS -o /dev/null -w '%{http_code}' "$short_url")"
assert_equals "404" "$status_code" "expired short link should return 404 after cleanup"

jq -n \
  --arg status "passed" \
  --arg short_url "$short_url" \
  --arg expires_at "$expired_at" \
  '{status:$status, short_url:$short_url, expires_at:$expires_at}' \
  > "$WORKSPACE/.eval/openai/verify_step2_result.json"
