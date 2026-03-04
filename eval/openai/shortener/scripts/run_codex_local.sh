#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-/work}"
SOURCE_FIXTURE="${SOURCE_FIXTURE:-$REPO_ROOT/eval/fixtures/draft/repo}"
RUN_ROOT="${RUN_ROOT:-/tmp/spec-agent-openai-run}"
WORKSPACE="${WORKSPACE:-$RUN_ROOT/workspace}"
MODEL="${CODEX_MODEL:-gpt-5-codex}"
REASONING_EFFORT="${CODEX_REASONING_EFFORT:-low}"
VERBOSE="${CODEX_VERBOSE:-0}"
CODEX_EXEC_MODE="${CODEX_EXEC_MODE:-unsafe}" # unsafe | full-auto
EXPORT_DIR="${EXPORT_DIR:-/out}"
DB_HOST="${DB_HOST:-postgres}"
RABBITMQ_HOST="${RABBITMQ_HOST:-localhost}"
APP_HOST="${APP_HOST:-http://localhost:8080}"
SPEC_AGENT_BIN="${SPEC_AGENT_BIN:-$RUN_ROOT/bin/spec-agent}"

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

build_spec_agent() {
  mkdir -p "$(dirname "$SPEC_AGENT_BIN")"
  (
    cd "$REPO_ROOT"
    GOCACHE="${GOCACHE:-$RUN_ROOT/.gocache}" go build -buildvcs=false -o "$SPEC_AGENT_BIN" ./cmd/spec-agent
  )
}

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
  local last_msg_file="$WORKSPACE/.eval/openai/artifacts/${name}_last_message.md"

  local mode_flag="--dangerously-bypass-approvals-and-sandbox"
  if [ "$CODEX_EXEC_MODE" = "full-auto" ]; then
    mode_flag="--full-auto"
  fi

  echo "running codex prompt: $name (model=$MODEL reasoning=$REASONING_EFFORT)"
  if [ "$VERBOSE" = "1" ]; then
    codex exec \
      --cd "$WORKSPACE" \
      --model "$MODEL" \
      "$mode_flag" \
      --skip-git-repo-check \
      --output-last-message "$last_msg_file" \
      -c "model_reasoning_effort=\"$REASONING_EFFORT\"" \
      - < "$prompt_file" | tee "$log_file"
    return
  fi

  if ! codex exec \
    --cd "$WORKSPACE" \
    --model "$MODEL" \
    "$mode_flag" \
    --skip-git-repo-check \
    --output-last-message "$last_msg_file" \
    -c "model_reasoning_effort=\"$REASONING_EFFORT\"" \
    - < "$prompt_file" >"$log_file" 2>&1; then
    echo "codex prompt failed: $name"
    echo "last log lines:"
    tail -n 80 "$log_file" || true
    return 1
  fi
  echo "codex prompt finished: $name"
}

export_artifacts() {
  local exit_code="$1"
  local status="failed"
  if [ "$exit_code" -eq 0 ]; then
    status="passed"
  fi

  local out_root="$WORKSPACE/.eval/openai"
  local out_artifacts="$out_root/artifacts"

  for f in mini_bugfix_report.json verify_mini_bugfix_result.json local_run_summary.json; do
    if [ -f "$out_root/$f" ]; then
      cp "$out_root/$f" "$DEST_DIR/$f"
    fi
  done

  if [ -d "$out_artifacts" ]; then
    mkdir -p "$DEST_DIR/artifacts"
    cp -R "$out_artifacts/." "$DEST_DIR/artifacts/"
  fi

  if [ -d "$WORKSPACE/.git" ]; then
    (
      cd "$WORKSPACE"
      git add -A
      git diff --binary --cached -- internal cmd migrations spec_changes > "$DEST_DIR/changes.patch" || true
      git diff --name-only --cached -- internal cmd migrations spec_changes > "$DEST_DIR/changed_files.txt" || true
    )
  else
    : > "$DEST_DIR/changes.patch"
    : > "$DEST_DIR/changed_files.txt"
  fi

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

echo "building fresh spec-agent binary..."
build_spec_agent

echo "injecting fresh prompts via spec-agent init..."
(
  cd "$WORKSPACE"
  "$SPEC_AGENT_BIN" init >/dev/null
)

echo "bootstrapping baseline specs for existing code..."
bash "$REPO_ROOT/eval/openai/shortener/scripts/bootstrap_existing_specs.sh"

echo "preparing eval-only runtime patch..."
bash "$REPO_ROOT/eval/openai/shortener/scripts/disable_rabbitmq_for_eval.sh"

echo "recording prepared baseline..."
(
  cd "$WORKSPACE"
  git add -A
  git commit -qm "prepared baseline with prompts/specs/eval patch"
)

run_prompt "mini_bugfix" "$REPO_ROOT/eval/openai/shortener/prompts/mini_bugfix.md"
bash "$REPO_ROOT/eval/openai/shortener/scripts/verify_mini_bugfix.sh"

jq -n \
  --arg status "passed" \
  --arg model "$MODEL" \
  --arg completed_at "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  '{status:$status, model:$model, completed_at:$completed_at}' \
  > "$WORKSPACE/.eval/openai/local_run_summary.json"

echo "local Codex eval completed successfully"
