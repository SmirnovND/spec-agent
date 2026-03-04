#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-/work}"
WORKSPACE="${WORKSPACE:-$REPO_ROOT/eval/fixtures/draft/repo}"
MODEL="${CODEX_MODEL:-gpt-5-codex}"

if ! command -v codex >/dev/null 2>&1; then
  echo "codex CLI is not installed"
  exit 1
fi

if ! codex login status >/dev/null 2>&1; then
  echo "codex login is required. Run 'codex login' on host and retry."
  exit 1
fi

mkdir -p "$WORKSPACE/.eval/openai"
export WORKSPACE
export SERVER_MODE="process"
export DB_HOST="${DB_HOST:-postgres}"
export RABBITMQ_HOST="${RABBITMQ_HOST:-rabbitmq}"
export APP_HOST="${APP_HOST:-http://localhost:8080}"

run_prompt() {
  local prompt_file="$1"
  codex exec \
    --cd "$WORKSPACE" \
    --model "$MODEL" \
    --full-auto \
    --skip-git-repo-check \
    - < "$prompt_file"
}

run_prompt "$REPO_ROOT/eval/openai/shortener/prompts/step1.md"
bash "$REPO_ROOT/eval/openai/shortener/scripts/verify_step1.sh"

run_prompt "$REPO_ROOT/eval/openai/shortener/prompts/step2.md"
bash "$REPO_ROOT/eval/openai/shortener/scripts/verify_step2.sh"

jq -n \
  --arg status "passed" \
  --arg model "$MODEL" \
  --arg completed_at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{status:$status, model:$model, completed_at:$completed_at}' \
  > "$WORKSPACE/.eval/openai/local_run_summary.json"

echo "local Codex eval completed successfully"
