#!/usr/bin/env bash
set -euo pipefail
: "${WORKSPACE:?WORKSPACE is required}"
: "${REPO_ROOT:?REPO_ROOT is required}"

SPEC_AGENT_BIN="${SPEC_AGENT_BIN:-$REPO_ROOT/bin/spec-agent}"
if [ ! -x "$SPEC_AGENT_BIN" ]; then
  mkdir -p "$REPO_ROOT/bin"
  GOCACHE="${GOCACHE:-$REPO_ROOT/.gocache}" go build -o "$SPEC_AGENT_BIN" "$REPO_ROOT/cmd/spec-agent"
fi

(
  cd "$WORKSPACE"
  "$SPEC_AGENT_BIN" init >/dev/null
)

mkdir -p "$WORKSPACE/spec_changes"
TARGET="$WORKSPACE/internal/router/router.go"

perl -0777 -i -pe 's/return nil/return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {\n\t\t\thttp.Error(w, "Service initialization failed", http.StatusInternalServerError)\n\t\t})/s' "$TARGET"

cat > "$WORKSPACE/internal/router/router.md" <<'MD'
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=internal/router/router -->
<!-- SPEC:KIND=other -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->

# Router

## Responsibility
Регистрирует HTTP маршруты и middleware для сервера.

## Inputs
- DI container с зависимостями приложения

## Outputs
- HTTP handler для обслуживания запросов

## Business Logic
1. Роутер должен регистрировать обязательные healthcheck маршруты.
2. При ошибке инициализации зависимостей должен возвращать fallback handler с кодом 500.

## Flow
1. Извлечь зависимости из DI контейнера.
2. При ошибке инициализации вернуть fallback handler.
3. Иначе зарегистрировать middleware и маршруты.

## Links
- uses: [HealthcheckController](../controllers/README.md#healthcheck)

## Dependencies
- [HealthcheckController](../controllers/README.md)

## Errors
- Ошибка инициализации зависимостей DI.
MD

cat > "$WORKSPACE/spec_changes/20260304_1510_bugfix_router_di_fallback.md" <<'MD'
# bugfix_router_di_fallback

- mode: bugfix
- root_cause: router returned nil handler when DI initialization failed
- fix: return explicit fallback HTTP handler with 500 status
- spec_update: document fallback handler behavior in router spec
MD
