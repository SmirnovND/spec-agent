#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

mkdir -p "$WORKSPACE/.eval/openai"

TARGET_GO="$WORKSPACE/internal/router/router.go"
TARGET_SPEC="$WORKSPACE/internal/router/router.md"
TARGET_TEST="$WORKSPACE/internal/router/router_fallback_test.go"

if ! grep -q "http.StatusInternalServerError" "$TARGET_GO"; then
  echo "missing http.StatusInternalServerError in router.go"
  exit 1
fi

if ! grep -Eqi "fallbackHandler|return .*fallback|return http\\.HandlerFunc|func fallback" "$TARGET_GO"; then
  echo "missing fallback handler implementation/return in router.go"
  exit 1
fi

if [ ! -f "$TARGET_TEST" ]; then
  echo "missing router fallback test file"
  exit 1
fi

if ! grep -q "TestHandlerFallbackOnDIError" "$TARGET_TEST"; then
  echo "missing required test name in fallback test file"
  exit 1
fi

if ! grep -Eqi "fallback|резервн|фолбэк|fallback-обработчик|fallback обработчик" "$TARGET_SPEC"; then
  echo "router spec is not updated with fallback behavior"
  exit 1
fi

if ! grep -q "500" "$TARGET_SPEC"; then
  echo "router spec is not updated with HTTP 500 fallback behavior"
  exit 1
fi

if ! find "$WORKSPACE/spec_changes" -maxdepth 1 -type f -name '*bugfix*.md' | grep -q .; then
  echo "missing spec_changes bugfix file"
  exit 1
fi

(
  cd "$WORKSPACE"
  GOCACHE="${GOCACHE:-$WORKSPACE/.gocache}" go test ./internal/router -run TestHandlerFallbackOnDIError -count=1
)

jq -n \
  --arg status "passed" \
  '{status:$status}' \
  > "$WORKSPACE/.eval/openai/verify_mini_bugfix_result.json"
