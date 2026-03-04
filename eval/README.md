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

## Draft Fixture (optional)

To sync external fixture repo (`SmirnovND/draft`):

1. Configure source in `eval/fixtures/draft/source.yaml`
2. Run:

```bash
bash eval/scripts/fetch_draft_fixture.sh
```

For stable runs, use pinned commit SHA in `ref` (already configured in `source.yaml`).
