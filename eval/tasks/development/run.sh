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

if ! grep -q 'r.Get("/healthz", healthcheckController.HandlePing)' "$TARGET"; then
  perl -0777 -i -pe 's/r\.Get\("\/ping", healthcheckController\.HandlePing\)\n/r.Get("\/ping", healthcheckController.HandlePing)\n\tr.Get("\/healthz", healthcheckController.HandlePing)\n/s' "$TARGET"
fi

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
2. Роутер должен поддерживать маршруты `/ping` и `/healthz` для healthcheck.

## Flow
1. Извлечь зависимости из DI контейнера.
2. Создать роутер и зарегистрировать middleware.
3. Зарегистрировать маршруты `/ping` и `/healthz` на healthcheck handler.

## Links
- uses: [HealthcheckController](../controllers/README.md#healthcheck)

## Dependencies
- [HealthcheckController](../controllers/README.md)

## Errors
- Ошибка инициализации зависимостей DI.
MD

cat > "$WORKSPACE/spec_changes/20260304_1520_development_add_healthz_route.md" <<'MD'
# development_add_healthz_route

- mode: development
- step_1: keep existing /ping route
- step_2: add /healthz route to the same healthcheck handler
- step_3: update router specification and links
MD
