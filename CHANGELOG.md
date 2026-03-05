# Changelog

All notable changes to this project are documented in this file.

## v2.0.3 - 2026-03-05

### Features

- Added `spec-agent version` command to print current CLI version.
- Added build-time version injection support via `-ldflags "-X main.version=..."`.

### Documentation

- Updated installation verification to use `spec-agent version`.
- Added example of versioned build command in README.

## v2.0.2 - 2026-03-05

### Fixes

- Fixed Go module major versioning for `v2+`:
  - changed module path to `github.com/SmirnovND/spec-agent/v2`
  - updated all internal imports to use `/v2/...`

### Documentation

- Updated installation command in `README.md` to:
  - `go install github.com/SmirnovND/spec-agent/v2/cmd/spec-agent@latest`

## v2.0.1 - 2026-03-04

### Documentation

- Synchronized main `README.md` with actual v2 architecture:
  - corrected project tree paths to `internal/fs/assets/shared/*`
  - updated specification format section to current `spec_rules`
  - refreshed usecase example to match marker/sections/links requirements
  - removed outdated non-canonical service example

## v2.0.0 - 2026-03-04

### Breaking Changes

- Prompt architecture switched to a single entrypoint model:
  - `.spec_agent/prompts/entrypoint.md`
  - `.spec_agent/prompts/core/*`
  - `.spec_agent/prompts/modes/*`
- Platform-specific adapter profiles were replaced by task-mode workflows:
  - `bugfix`, `development`, `refactor`, `tests`, `analysis`
- Shared prompt language policy changed:
  - `core` and `modes` prompts are in English
  - specification content/output remains Russian

### Features

- Added machine-readable prompt source layer (`prompt.yaml`) for shared prompts.
- Added prompt generator CLI:
  - `go run ./cmd/promptgen`
  - `go run ./cmd/promptgen -check`
- Added single routing entry prompt (`entrypoint`) with mandatory task-mode selection.
- Added task-mode selector and dedicated mode workflows.
- Added `internal/promptgen` package with generation and validation logic.
- Added tests for prompt generation checks.

### Behavior Changes

- Re-running `spec-agent init` or `spec-agent init-php` no longer overwrites existing files.
- `init` now creates only missing files for config/prompts/examples.

### Documentation

- Updated README and assets docs for:
  - entrypoint-based prompt routing
  - task modes
  - machine-readable prompt generation
  - language policy
