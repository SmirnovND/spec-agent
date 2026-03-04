# Eval MVP

This folder contains a minimal prompt-evaluation framework for `spec-agent`.

## Metrics

Each task run computes:
- `completeness`: passed checkpoints / total checkpoints
- `scope_violations`: changed files outside `allowed_paths`
- `test_pass_rate`: passed test commands / total test commands
- `spec_rule_violations`: violations in changed spec files (`*.md` with `SPEC:FILE=true`)

Gating thresholds are stored in:
- `eval/baselines/thresholds.yaml`

## Structure

- `eval/tasks/<task_id>/task.md` - task description
- `eval/tasks/<task_id>/config.yaml` - fixture, allowed paths, commands
- `eval/tasks/<task_id>/checks.yaml` - completeness checkpoints
- `eval/fixtures/<fixture>/` - fixture source project
- `eval/results/*.json` - run outputs
- `cmd/eval-runner/main.go` - runner CLI

## Local Run

Before running eval on `draft`, sync fixture:

```bash
bash eval/scripts/fetch_draft_fixture.sh
```

Run all tasks:

```bash
go run ./cmd/eval-runner
```

Run selected tasks:

```bash
go run ./cmd/eval-runner -tasks bootstrap,bugfix,development
```

Output JSON is written to `eval/results/<timestamp>.json`.

## CI

Workflow:
- `.github/workflows/eval.yml`

Behavior:
- quick eval on PR/push for relevant paths
- full eval on schedule / manual dispatch
- JSON outputs uploaded as artifacts
- jobs run in short-lived container (`golang:1.24`) and are discarded after completion

## OpenAI Agent Eval (Mini Bugfix)

Workflow:
- `.github/workflows/eval-openai-shortener.yml`

What it does:
- syncs pinned `draft` fixture
- injects fresh prompts via `spec-agent init`
- bootstraps baseline specs for existing modules
- runs OpenAI coding agent on one small task (`mini_bugfix`)
- verifies behavior, scope and required artifacts
- uploads run artifacts from `eval/fixtures/draft/repo/.eval/openai/`

Requirements:
- GitHub repository secret `OPENAI_API_KEY` must be configured
- run workflow manually via `workflow_dispatch`

Prompt files:
- `eval/openai/shortener/prompts/mini_bugfix.md`

Verification scripts:
- `eval/openai/shortener/scripts/verify_mini_bugfix.sh`

Local isolated run (Codex CLI in container):
- run `codex login` once on host machine
- `make eval-openai-shortener-local`

Optional runtime knobs:
- `CODEX_VERBOSE=1` - stream full Codex output to console (default `0`, quiet mode)
- `CODEX_REASONING_EFFORT=low|medium|high` - speed/quality tradeoff (default `low`)
- `CODEX_MODEL=...` - override model
- `CODEX_EXEC_MODE=unsafe|full-auto` - Codex execution mode (default `unsafe`)

This target:
- syncs pinned `draft` fixture
- starts `postgres` + `runner` via `docker compose`
- builds fresh `spec-agent` binary from current repository state
- runs `spec-agent init` inside ephemeral workspace to inject latest prompt templates
- bootstraps baseline specs for existing modules (router/healthcheck layers)
- applies eval-only patch to disable mandatory RabbitMQ startup in `cmd/server/main.go`
- runs `codex exec` for `mini_bugfix`
- runs `verify_mini_bugfix` after agent completion
- writes artifacts to `eval/results/openai-shortener/<run_id>/`

Auth mode:
- local runner uses mounted host profile `${HOME}/.codex -> /root/.codex`
- no `OPENAI_API_KEY` export required for local run

Artifact export policy:
- only compact artifacts are exported (`summary.json`, mini report, verify result, codex logs)
- only allowed code diff is exported as `changes.patch` for paths:
  - `internal/**`
  - `cmd/**`
  - `migrations/**`
  - `spec_changes/**`
- full workspace is temporary and removed with container lifecycle

## Draft Fixture (optional)

To sync external fixture repo (`SmirnovND/draft`):

1. Configure source in `eval/fixtures/draft/source.yaml`
2. Run:

```bash
bash eval/scripts/fetch_draft_fixture.sh
```

For stable runs, use pinned commit SHA in `ref` (already configured in `source.yaml`).
