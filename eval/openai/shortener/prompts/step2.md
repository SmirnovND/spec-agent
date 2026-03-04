Continue from current workspace state.

Task (Russian source):
"Доработай функционал создания короткой ссылки. Необходимо принимать на вход дату окончания жизни
ссылки. так же написать консольную команду которая будет удалять из базы ссылки с прошедшей датой
Применяем миграции make migrate-up".

Hard requirements:
1. Extend `POST /api/v1/short-links` request with optional `expires_at` field in RFC3339.
2. Persist expiration in DB.
3. Add console command at `cmd/crons/cleanup_short_links/main.go`:
   - accepts config path arg, same style as existing commands
   - deletes expired links from DB
4. Apply schema updates via `make migrate-up`.
5. Keep `GET /{hash}` redirect behavior:
   - active link -> HTTP 302
   - expired or missing link -> HTTP 404

Repository constraints:
- Keep existing architecture.
- Keep backward compatibility for requests without `expires_at`.
- Keep code compilable.

Output contract:
- Create file `.eval/openai/step2_report.json` with JSON:
  {
    "task": "shortener_step2",
    "status": "done",
    "changed_files": ["..."],
    "cleanup_command": "go run ./cmd/crons/cleanup_short_links/main.go ./cmd/server/config.yaml",
    "notes": "..."
  }
