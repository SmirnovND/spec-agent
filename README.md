# spec-agent

CLI инструмент для управления спецификациями, управляемыми архитектурой (spec-driven architecture). Позволяет определять структуру кода через спецификации, автоматизировать их проверку и визуализировать зависимости.

## Основной концепт

Спецификации — это MD-файлы рядом с кодом, которые:
- Определяют архитектурные решения и контексты
- Описывают контракты и ограничения компонентов
- Управляют зависимостями между компонентами
- Документируют причины и ответственности

Инструмент помогает проверить консистентность спецификаций и понять граф их зависимостей.

## Установка

```bash
go install github.com/SmirnovND/spec-agent/cmd/spec-agent@latest
```

Проверка установки:
```bash
spec-agent --help
```

## Быстрый старт

### 1. Инициализируйте проект

```bash
spec-agent init
```

Создаст `.spec_agent/config.yaml` с путями к спецификациям.
При повторном запуске `init` существующие файлы не перезаписываются: добавляются только отсутствующие.

Для PHP-проекта:

```bash
spec-agent init-php
```
При повторном запуске `init-php` существующие файлы не перезаписываются: добавляются только отсутствующие.

### 2. Обновите конфиг

```yaml
# .spec_agent/config.yaml
roots:
  - internal/services
  - internal/repositories
exclude:
  - cmd/staticlint
  - cmd/server
```

Укажите директории, где находятся корневые спецификации (MD-файлы).
Через `exclude` можно исключить технические папки из сканирования/экспорта.

### 3. Запустите сервер

```bash
spec-agent serve
```

Откройте в браузере: **http://localhost:8080**

### 4. Просмотрите граф

```bash
spec-agent graph
```

Увидите структуру зависимостей между спеками.

## Использование

### Инициализация проекта

```bash
spec-agent init
```

Создаёт структуру:
- `.spec_agent/config.yaml` — конфиг с корневыми путями для поиска спецификаций
- `spec_changes/` — директория для отслеживания изменений

Для PHP-проектов используйте отдельную команду:

```bash
spec-agent init-php
```

Дополнительно создаётся отдельный профиль:
- `.spec_agent/php/prompts/base/spec_rules.md` — правила для PHP спецификаций
- `.spec_agent/php/prompts/base/agent_prompt.md` — базовый prompt для PHP-агента
- `.spec_agent/php/prompts/base/workflow.md` — workflow для PHP
- `.spec_agent/php/examples/` — примеры спецификаций для PHP

Архитектурное правило слоёв (зафиксировано в промптах):
- Сценарная бизнес-оркестрация должна находиться в `usecase`.
- `service` реализует отдельные шаги сценария и технические операции, вызываемые из `usecase`.
- Логика полноценного пользовательского/бизнес-сценария не должна размещаться в `service`.

Для интеграции с Zenflow можно сразу создать кастомный workflow:

```bash
spec-agent init --zenflow
```

Отличие источников промтов:
- `init` берёт промты из `internal/fs/assets/go/prompts/base`
- `init --zenflow` добавляет zenflow-набор из `internal/fs/assets/go/prompts/zenflow`
- `init` также добавляет общий слой в `.spec_agent/prompts/`:
  - `entrypoint.md`: единый вход и порядок чтения инструкций
  - `core/*`: единый контракт, базовый workflow и task-mode selector
  - `modes/*`: режимы по типу задачи (`bugfix`, `development`, `refactor`, `tests`, `analysis`)

Языковая политика промптов:
- `.spec_agent/prompts/core/*` и `.spec_agent/prompts/modes/*` — на английском.
- Спецификации рядом с кодом (`*.md`) и бизнес-текст — на русском.

Дополнительно создаётся файл:
- `.zenflow/workflows/spec-agent-spec-driven.md` — custom workflow с шагами `Planning`, `Technical Specification`, `Specification Review`, `Implementation`, `Review & Wrap-Up`
- В workflow используются артефакты Zenflow: `{@artifacts_path}/plan.md`, `{@artifacts_path}/spec.md`, `{@artifacts_path}/report.md`
- Для каждого шага используются отдельные обязательные инструкции:
  - `.spec_agent/prompts/zenflow/planning.md`
  - `.spec_agent/prompts/zenflow/technical_specification.md`
  - `.spec_agent/prompts/zenflow/specification_review.md`
  - `.spec_agent/prompts/zenflow/implementation.md`
  - `.spec_agent/prompts/zenflow/review_wrap_up.md`

Для PHP-профиля:

```bash
spec-agent init-php --zenflow
```

Создаётся отдельный workflow:
- `.zenflow/workflows/spec-agent-php-spec-driven.md`
- шаги ссылаются на `.spec_agent/php/prompts/zenflow/*.md`
- общий слой `.spec_agent/prompts/entrypoint.md`, `.spec_agent/prompts/core` и `.spec_agent/prompts/modes` также создаётся

### Просмотр спецификаций в браузере

**Способ 1: Встроенный веб-сервер (рекомендуется)**

```bash
spec-agent serve -p 8080
```

Запускает встроенный HTTP-сервер на `http://localhost:8080`:
- Автоматически генерирует HTML если спеки ещё не экспортированы
- Обслуживает статические файлы из `.spec_agent/build/`
- Прекращает работу по `Ctrl+C`
- Не требует nginx, node или docker

Опции:
```bash
spec-agent serve -p 3000              # Другой порт
spec-agent serve --host 0.0.0.0       # Доступен для других хостов
```

**Способ 2: Экспорт в статичный HTML**

```bash
spec-agent export
```

Генерирует HTML файлы в `.spec_agent/build/`:
- `index.html` — главная страница с оглавлением
- `{spec_name}.md.html` — отдельные страницы спеков
- Можно открыть как локальный файл в браузере

### Просмотр графа зависимостей

```bash
spec-agent graph
```

Анализирует спецификации и выводит информацию:
- Найденные root-спеки (на которые никто не ссылается)
- Граф зависимостей (кол-во узлов и рёбер)
- Определяет структуру и взаимосвязи

## Структура проекта

```
spec-agent/
├── cmd/spec-agent/
│   └── main.go               # Точка входа приложения
├── cmd/promptgen/
│   └── main.go               # Генерация markdown из prompt.yaml
├── internal/
│   ├── cli/                  # Команды CLI (Cobra)
│   │   ├── root.go           # Корневая команда
│   │   ├── init.go           # spec-agent init
│   │   ├── graph.go          # spec-agent graph
│   │   ├── export.go         # spec-agent export
│   │   └── serve.go          # spec-agent serve
│   ├── spec/                 # Логика работы со спецификациями
│   │   ├── model.go          # Структуры: Spec, Graph, Node, Edge
│   │   ├── parser.go         # Парсинг MD-файлов
│   │   ├── graph.go          # Построение графа зависимостей
│   │   └── exporter.go       # Генерация HTML
│   ├── config/
│   │   └── config.go         # Загрузка .spec_agent/config.yaml
│   └── fs/
│       └── init.go           # Инициализация проекта
├── assets/
│   ├── examples/             # Примеры спецификаций
│   ├── shared/prompts/       # Общий слой: entrypoint + core + task modes
│   ├── shared/prompt_specs/  # Машинно-читабельные prompt.yaml
│   ├── go/prompts/           # Go-профиль: base + zenflow
│   └── php/prompts/          # PHP-профиль: base + zenflow
├── go.mod
├── go.sum
└── README.md
```

## Зависимости

- **[Cobra](https://github.com/spf13/cobra)** — фреймворк для CLI команд
- **[gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)** — парсинг YAML конфига

## Формат спецификаций

Спецификации — это Markdown-файлы со следующей структурой:

```markdown
# Название компонента

## Контекст
Описание задачи, роли в системе, исторический контекст.

## Ответственность
Что этот компонент отвечает за реализацию.

## Контракты
Интерфейсы, структуры данных, API.

## Ограничения
Constraint'ы, assumptions, limitations.

## Зависимости
- Ссылка на другую спецификацию: [Component Name](../other/spec.md)
```

## Конфигурация

`.spec_agent/config.yaml`:

```yaml
roots:
  - internal/controllers  # Где искать спецификации
  - internal/middleware
exclude:
  - cmd/staticlint        # Что исключить из сканирования
  - cmd/server
```

## Примеры

### Пример спецификации usecase

```markdown
# CreateUserUseCase

## Responsibility
Реализует бизнес-процесс создания нового пользователя со всеми необходимыми проверками.

## Inputs
- Email пользователя
- Пароль
- Имя пользователя

## Outputs
- Созданный объект User с ID

## Business Rules
1. Email должен быть уникальным
2. Пароль должен быть минимум 8 символов
3. Новый пользователь создается в неактивном состоянии

## Flow
1. Валидирует входные данные
2. Проверяет уникальность email → calls: [UserRepository](../repositories/user_repository.md)
3. Хеширует пароль → calls: [CryptoService](../services/crypto_service.md)
4. Создает запись в БД → writes: [UserRepository](../repositories/user_repository.md)
5. Отправляет письмо подтверждения → calls: [EmailService](../services/email_service.md)

## Dependencies
- [UserRepository](../repositories/user_repository.md)
- [CryptoService](../services/crypto_service.md)
- [EmailService](../services/email_service.md)

## Errors
- ErrEmailExists — email уже зарегистрирован
- ErrInvalidEmail — некорректный формат
- ErrWeakPassword — слабый пароль
```

### Пример спецификации сервиса

```markdown
# UserRepository

## Responsibility
Управляет доступом к данным пользователей в базе данных.

## Inputs
- User object for create/update
- User ID for read/delete
- Filter parameters for list

## Outputs
- User object(s) or error

## Dependencies
- [Database Driver](../db/connection.md)

## Contract
- Create(ctx, user) → (userWithID, error)
- GetByID(ctx, id) → (user, error)
- GetByEmail(ctx, email) → (user, error)
- Update(ctx, user) → error
- Delete(ctx, id) → error
```

## Разработка

### Сборка локально

```bash
go build -o spec-agent ./cmd/spec-agent
./spec-agent --help
```

### Тестирование

```bash
go test ./...
```

### Генерация prompt-маркировки

Исходники промптов хранятся в `internal/fs/assets/shared/prompt_specs/**/prompt.yaml`.
Сгенерировать markdown:

```bash
go run ./cmd/promptgen
```

Проверить, что markdown актуален:

```bash
go run ./cmd/promptgen -check
```

### Структура команд

Каждая команда в `internal/cli/` соответствует CLI команде:
- `root.go` — корневая команда и регистрация подкоманд
- `export.go` — генерация HTML
- `serve.go` — встроенный веб-сервер
- `graph.go` — анализ зависимостей
- `init.go` — инициализация проекта

### Логика обработки спецификаций

В `internal/spec/` находится основная логика:
- `parser.go` — парсинг MD-файлов в структуру Spec
- `graph.go` — построение графа зависимостей
- `exporter.go` — генерация HTML с навигацией
