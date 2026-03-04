You are working on a Go backend in this workspace.

Нужно добавить в сервис базовый функционал сокращения ссылок.

Что ожидаю по результату:
- Появляется `POST /api/v1/short-links`, который принимает JSON с полем `url`.
- Внутри логика такая: берем входной URL, выделяем URN (путь + query), хэшируем URN и строим короткую ссылку как `<shortener_host>/<hash>`.
- `shortener_host` должен быть в конфиге как отдельная переменная.
- Сохраняем в PostgreSQL связку оригинального URL и короткого URL.
- В ответе на создание вернуть JSON вида:
  - `original_url`
  - `short_url`
- Появляется `GET /{hash}`, который делает redirect (302) на полный URL.
- `/ping` не ломаем.
- Если нужны миграции, подготовь их и примени через `make migrate-up`.

Работай в рамках текущей архитектуры проекта (controller/service/repository/container/router), без радикальной перестройки.

После выполнения сохрани короткий отчет в `.eval/openai/step1_report.json`:
{
  "task": "shortener_step1",
  "status": "done",
  "changed_files": ["..."],
  "notes": "..."
}
