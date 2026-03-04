#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

mkdir -p \
  "$WORKSPACE/internal/router" \
  "$WORKSPACE/internal/controllers" \
  "$WORKSPACE/internal/services" \
  "$WORKSPACE/internal/repositories" \
  "$WORKSPACE/spec_changes"

cat > "$WORKSPACE/internal/router/router.md" <<'MD'
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=internal/router/router -->
<!-- SPEC:KIND=other -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->

# Router

## Responsibility
Регистрирует HTTP-маршруты и базовые middleware приложения.

## Inputs
- DI container с зависимостями приложения.

## Outputs
- `http.Handler`, обслуживающий входящие HTTP-запросы.

## Business Logic
1. Маршрут `/ping` должен быть доступен всегда.
2. Роутер должен корректно обрабатывать `404` и `405`.

## Flow
1. Получить зависимости контроллеров через DI.
2. Создать роутер и подключить middleware.
3. Зарегистрировать endpoint'ы и служебные обработчики ошибок.

## Links
- uses: [Healthcheck Controller](../controllers/healthcheck_controller.md#internalcontrollershealthcheck_controller)

## Dependencies
- [Healthcheck Controller](../controllers/healthcheck_controller.md)

## Errors
- Ошибка DI при инициализации контроллеров.
MD

cat > "$WORKSPACE/internal/controllers/healthcheck_controller.md" <<'MD'
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=internal/controllers/healthcheck_controller -->
<!-- SPEC:KIND=controller -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->

# Healthcheck Controller

## Responsibility
Обрабатывает HTTP-запрос `GET /ping` и возвращает статус здоровья сервиса.

## Inputs
- HTTP-запрос.
- `HealthcheckService` из DI.

## Outputs
- JSON-ответ со статусом сервиса.

## Business Logic
1. Контроллер делегирует проверку состояния сервисному слою.
2. При ошибке проверки возвращает `500` и текст ошибки.

## Flow
1. Принять запрос `/ping`.
2. Вызвать `HealthcheckService.Check(ctx)`.
3. Вернуть `200` или `500` в зависимости от результата.

## Links
- calls: [Healthcheck Service](../services/healthcheck_service.md#internalserviceshealthcheck_service)

## Dependencies
- [Healthcheck Service](../services/healthcheck_service.md)

## Errors
- Ошибка проверки подключения к базе данных.
MD

cat > "$WORKSPACE/internal/services/healthcheck_service.md" <<'MD'
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=internal/services/healthcheck_service -->
<!-- SPEC:KIND=service -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->

# Healthcheck Service

## Responsibility
Проверяет доступность инфраструктурных зависимостей приложения.

## Inputs
- Контекст выполнения запроса.
- `HealthcheckRepository`.

## Outputs
- Структура статуса сервиса (`status=ok`) или ошибка.

## Business Logic
1. Проверка здоровья выполняется через репозиторий.
2. При успешной проверке возвращается `status=ok`.

## Flow
1. Вызвать `HealthcheckRepository.Ping(ctx)`.
2. Если ошибка, вернуть ее выше.
3. Иначе вернуть успешный статус.

## Links
- calls: [Healthcheck Repository](../repositories/healthcheck_repository.md#internalrepositorieshealthcheck_repository)

## Dependencies
- [Healthcheck Repository](../repositories/healthcheck_repository.md)

## Errors
- Ошибка подключения к БД.
MD

cat > "$WORKSPACE/internal/repositories/healthcheck_repository.md" <<'MD'
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=internal/repositories/healthcheck_repository -->
<!-- SPEC:KIND=repository -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->

# Healthcheck Repository

## Responsibility
Проверяет соединение с базой данных через текущий DB-драйвер.

## Inputs
- Контекст выполнения.
- Экземпляр `sqlx.DB`.

## Outputs
- `nil` при успешном ping.
- Ошибка при недоступности БД.

## Business Logic
1. Проверка делается через `db.PingContext`.
2. Репозиторий не содержит бизнес-логики.

## Flow
1. Получить запрос на ping.
2. Выполнить `PingContext`.
3. Вернуть результат выше по стеку.

## Links
- validates: [Healthcheck Service](../services/healthcheck_service.md#internalserviceshealthcheck_service)

## Dependencies
- [Healthcheck Service](../services/healthcheck_service.md)

## Errors
- Ошибка соединения с БД или таймаут.
MD

cat > "$WORKSPACE/spec_changes/20260304_0000_bootstrap_existing_specs.md" <<'MD'
# bootstrap_existing_specs

- mode: analysis
- action: create baseline specs for existing healthcheck/router modules
- goal: ensure repository starts in spec-ready state before task execution
MD
