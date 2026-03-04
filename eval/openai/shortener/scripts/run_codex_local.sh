#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-/work}"
SOURCE_FIXTURE="${SOURCE_FIXTURE:-$REPO_ROOT/eval/fixtures/draft/repo}"
RUN_ROOT="${RUN_ROOT:-/tmp/spec-agent-openai-run}"
WORKSPACE="${WORKSPACE:-$RUN_ROOT/workspace}"
MODEL="${CODEX_MODEL:-gpt-5-codex}"
EXPORT_DIR="${EXPORT_DIR:-/out}"
DB_HOST="${DB_HOST:-postgres}"
RABBITMQ_HOST="${RABBITMQ_HOST:-rabbitmq}"
APP_HOST="${APP_HOST:-http://localhost:8080}"

if ! command -v codex >/dev/null 2>&1; then
  echo "codex CLI is not installed"
  exit 1
fi

if ! codex login status >/dev/null 2>&1; then
  echo "codex login is required. Run 'codex login' on host and retry."
  exit 1
fi

if [ ! -d "$SOURCE_FIXTURE" ]; then
  echo "source fixture not found: $SOURCE_FIXTURE"
  exit 1
fi

RUN_ID="$(date -u +%Y%m%d_%H%M%S)"
DEST_DIR="$EXPORT_DIR/$RUN_ID"
mkdir -p "$RUN_ROOT" "$DEST_DIR"

prepare_workspace() {
  rm -rf "$WORKSPACE"
  mkdir -p "$WORKSPACE"
  rsync -a --delete \
    --exclude '.git' \
    "$SOURCE_FIXTURE/" "$WORKSPACE/"

  mkdir -p "$WORKSPACE/.eval/openai/artifacts"

  (
    cd "$WORKSPACE"
    git init -q
    git config user.email "eval@local"
    git config user.name "eval"
    git add -A
    git commit -qm "baseline"
  )
}

run_prompt() {
  local name="$1"
  local prompt_file="$2"
  local log_file="$WORKSPACE/.eval/openai/artifacts/${name}_codex.log"

  codex exec \
    --cd "$WORKSPACE" \
    --model "$MODEL" \
    --full-auto \
    --dangerously-bypass-approvals-and-sandbox \
    --skip-git-repo-check \
    - < "$prompt_file" | tee "$log_file"
}

export_artifacts() {
  local exit_code="$1"
  local status="failed"
  if [ "$exit_code" -eq 0 ]; then
    status="passed"
  fi

  local out_root="$WORKSPACE/.eval/openai"
  local out_artifacts="$out_root/artifacts"

  for f in step1_report.json step2_report.json verify_step1_result.json verify_step2_result.json local_run_summary.json; do
    if [ -f "$out_root/$f" ]; then
      cp "$out_root/$f" "$DEST_DIR/$f"
    fi
  done

  if [ -d "$out_artifacts" ]; then
    mkdir -p "$DEST_DIR/artifacts"
    cp -R "$out_artifacts/." "$DEST_DIR/artifacts/"
  fi

  (
    cd "$WORKSPACE"
    git add -A
    git diff --binary --cached -- internal cmd migrations spec_changes > "$DEST_DIR/changes.patch" || true
    git diff --name-only --cached -- internal cmd migrations spec_changes > "$DEST_DIR/changed_files.txt" || true
  )

  jq -n \
    --arg status "$status" \
    --arg model "$MODEL" \
    --arg run_id "$RUN_ID" \
    --arg workspace "$WORKSPACE" \
    --arg source_fixture "$SOURCE_FIXTURE" \
    --arg completed_at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    --argjson exit_code "$exit_code" \
    '{status:$status, model:$model, run_id:$run_id, source_fixture:$source_fixture, workspace:$workspace, completed_at:$completed_at, exit_code:$exit_code}' \
    > "$DEST_DIR/summary.json"

  echo "artifacts exported to: $DEST_DIR"
}

cleanup() {
  local code="$?"
  export_artifacts "$code"
}
trap cleanup EXIT

prepare_workspace

export WORKSPACE
export SERVER_MODE="process"
export DB_HOST
export RABBITMQ_HOST
export APP_HOST

run_prompt "step1" "$REPO_ROOT/eval/openai/shortener/prompts/step1.md"
bash "$REPO_ROOT/eval/openai/shortener/scripts/verify_step1.sh"

run_prompt "step2" "$REPO_ROOT/eval/openai/shortener/prompts/step2.md"
bash "$REPO_ROOT/eval/openai/shortener/scripts/verify_step2.sh"

jq -n \
  --arg status "passed" \
  --arg model "$MODEL" \
  --arg completed_at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{status:$status, model:$model, completed_at:$completed_at}' \
  > "$WORKSPACE/.eval/openai/local_run_summary.json"

echo "local Codex eval completed successfully"
