#!/usr/bin/env bash
set -euo pipefail
: "${WORKSPACE:?WORKSPACE is required}"
: "${REPO_ROOT:?REPO_ROOT is required}"

SPEC_AGENT_BIN="${SPEC_AGENT_BIN:-$REPO_ROOT/bin/spec-agent}"
if [ ! -x "$SPEC_AGENT_BIN" ]; then
  mkdir -p "$REPO_ROOT/bin"
  GOCACHE="${GOCACHE:-$REPO_ROOT/.gocache}" go build -buildvcs=false -o "$SPEC_AGENT_BIN" "$REPO_ROOT/cmd/spec-agent"
fi

(
  cd "$WORKSPACE"
  "$SPEC_AGENT_BIN" init >/dev/null
)

mkdir -p "$WORKSPACE/internal/router" "$WORKSPACE/spec_changes"

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
2. При ошибке инициализации зависимостей должен возвращать безопасный fallback handler.

## Flow
1. Извлечь зависимости из DI контейнера.
2. Создать роутер и зарегистрировать middleware.
3. Зарегистрировать healthcheck endpoint и служебные обработчики ошибок.

## Links
- uses: [HealthcheckController](../controllers/README.md#healthcheck)

## Dependencies
- [HealthcheckController](../controllers/README.md)

## Errors
- Ошибка инициализации зависимостей DI.
MD

cat > "$WORKSPACE/spec_changes/20260304_1500_bootstrap_router_spec.md" <<'MD'
# bootstrap_router_spec

- mode: development
- step_1: run spec-agent init in draft fixture
- step_2: create initial router specification for existing router code
MD
