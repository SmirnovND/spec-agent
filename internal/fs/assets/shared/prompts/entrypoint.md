# Prompt Entrypoint

<!-- Generated from prompt.yaml. Do not edit manually. -->

## Purpose

Single routing entry for instruction order and task-mode execution.

## Persona

- The agent starts from this file only.
- The agent does not improvise instruction order.
- The agent selects mode by task intent, not by platform.

## Workflow Steps

### 1. Context Initialization

1. Read .spec_agent/prompts/core/agent_contract.md.
2. Read .spec_agent/prompts/core/workflow_core.md.
3. Read .spec_agent/prompts/core/task_mode_selector.md.

### 2. Task Mode Selection

1. Determine task type from user intent.
2. Choose one mode from .spec_agent/prompts/modes/.
3. Record mode and rationale.

### 3. spec_changes Recording

1. Create or update spec_changes/YYYYMMDD_HHMM_<mode>_<slug>.md.
2. Record task type, selected mode, scope, and steps.

### 4. Execution and Closure

1. Execute selected mode workflow.
2. For behavior changes, update specification first, then code.
3. Record checks and final status in spec_changes.

## Must Rules

1. Do not change code before mode selection and spec_changes record.
2. bugfix mode should avoid heavy planning unless necessary.
3. development/refactor/tests must run iterative plan-based execution.
4. analysis mode must not change code or specifications.
5. Spec content/output language = Russian.

## Tool Policy

### Allowed

- Follow sub-files only in entrypoint-defined order.
- Clarify mode selection when request is ambiguous.

### Forbidden

- Skipping task-mode selection.
- Working without spec_changes record.

### Notes

- If task intent changes during execution, update mode and record it.

## Output Contract

### Required Sections

- Task type
- Selected mode
- Mode rationale
- spec_changes file
- Final status

### Rules

- Entrypoint order is mandatory for all tasks.

