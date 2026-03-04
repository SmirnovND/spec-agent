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

## OpenAI Agent Eval (Shortener Scenario)

Workflow:
- `.github/workflows/eval-openai-shortener.yml`

What it does:
- syncs pinned `draft` fixture
- runs OpenAI coding agent (via `openai/codex-action`) with task #1 prompt
- verifies API behavior by running service in Docker container and calling HTTP endpoints
- runs OpenAI coding agent with task #2 prompt
- verifies expiration + cleanup-cron behavior
- uploads run artifacts from `eval/fixtures/draft/repo/.eval/openai/`

Requirements:
- GitHub repository secret `OPENAI` must be configured
- run workflow manually via `workflow_dispatch`

Prompt files:
- `eval/openai/shortener/prompts/step1.md`
- `eval/openai/shortener/prompts/step2.md`

Verification scripts:
- `eval/openai/shortener/scripts/verify_step1.sh`
- `eval/openai/shortener/scripts/verify_step2.sh`

Local isolated run (Codex CLI in container):
- run `codex login` once on host machine
- `make eval-openai-shortener-local`

Optional runtime knobs:
- `CODEX_VERBOSE=1` - stream full Codex output to console (default `0`, quiet mode)
- `CODEX_REASONING_EFFORT=low|medium|high` - speed/quality tradeoff (default `low`)
- `CODEX_MODEL=...` - override model

This target:
- syncs pinned `draft` fixture
- starts `postgres` + `rabbitmq` + `runner` via `docker compose`
- runs `codex exec` for step1 and step2 in an ephemeral workspace inside `runner`
- executes runtime verification scripts after each step
- writes artifacts to `eval/results/openai-shortener/<run_id>/`

Auth mode:
- local runner uses mounted host profile `${HOME}/.codex -> /root/.codex`
- no `OPENAI_API_KEY` export required for local run

Artifact export policy:
- only compact artifacts are exported (`summary.json`, reports, verify results, codex logs)
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
