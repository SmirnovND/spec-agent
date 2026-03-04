#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/common.sh"

prepare_server_config
apply_migrations
build_server_image
start_server_container
trap 'stop_server_container' EXIT
wait_http_ready

mkdir -p "$WORKSPACE/.eval/openai"

CREATE_PAYLOAD="$WORKSPACE/.eval/openai/step1_create_payload.json"
CREATE_RESPONSE="$WORKSPACE/.eval/openai/step1_create_response.json"

cat > "$CREATE_PAYLOAD" <<'JSON'
{
  "url": "https://example.com/products/123?utm=abc"
}
JSON

curl -fsS \
  -X POST "http://localhost:8080/api/v1/short-links" \
  -H "Content-Type: application/json" \
  --data-binary "@$CREATE_PAYLOAD" \
  > "$CREATE_RESPONSE"

original_url="$(extract_json_field '.original_url' "$CREATE_RESPONSE")"
short_url="$(extract_json_field '.short_url' "$CREATE_RESPONSE")"

assert_equals "https://example.com/products/123?utm=abc" "$original_url" "original_url mismatch"
assert_prefix "http://localhost:8080/" "$short_url" "short_url host mismatch"

redirect_headers="$WORKSPACE/.eval/openai/step1_redirect_headers.txt"
status_code="$(curl -sS -o /dev/null -D "$redirect_headers" -w '%{http_code}' "$short_url")"
location_header="$(awk -F': ' 'tolower($1)=="location" {print $2}' "$redirect_headers" | tr -d '\r' | tail -n1)"

assert_equals "302" "$status_code" "redirect status mismatch"
assert_equals "https://example.com/products/123?utm=abc" "$location_header" "redirect location mismatch"

jq -n \
  --arg status "passed" \
  --arg original_url "$original_url" \
  --arg short_url "$short_url" \
  '{status:$status, original_url:$original_url, short_url:$short_url}' \
  > "$WORKSPACE/.eval/openai/verify_step1_result.json"
