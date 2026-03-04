# Workflow (Core)

<!-- Generated from prompt.yaml. Do not edit manually. -->

## Purpose

Base execution workflow for all agent systems.

## Persona

- The agent selects a task mode before execution.
- The agent does not apply one workflow to every task type.
- The agent leaves a traceable record in spec_changes.

## Workflow Steps

### 1. Start

1. Always start from .spec_agent/prompts/entrypoint.md.

### 2. Mode Selection

1. Apply task_mode_selector.
2. Select one mode: bugfix, development, refactor, tests, analysis.

### 3. Context Analysis

1. Identify entry-point component.
2. Read related specifications.
3. Define scope boundaries.

### 4. spec_changes Record

1. Create spec_changes/YYYYMMDD_HHMM_<task_type>_<slug>.md.
2. Record selected mode, scope, and planned steps.

### 5. Mode Execution

1. Follow selected mode workflow.
2. For behavior changes, sync specification first, then code.

### 6. Verification and Closure

1. Run checks required by selected mode.
2. Update spec_changes with outcomes and final status.

## Must Rules

1. Selected mode must be explicit in spec_changes.
2. Do not use heavy development workflow for a simple bugfix without reason.
3. Do not skip final record update in spec_changes.

## Tool Policy

### Allowed

- Use only steps relevant to selected task mode.
- Run local checks required by mode.

### Forbidden

- Mixing multiple modes without explicit justification.

### Notes

- If uncertainty is high, record assumptions in spec_changes.

## Output Contract

### Required Sections

- Task type and selected mode
- Mode steps/plan
- Implemented changes
- Verification
- Risks

### Rules

- Every step in spec_changes must have status and outcome.

