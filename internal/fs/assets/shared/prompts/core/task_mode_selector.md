# Task Mode Selector (Core)

<!-- Generated from prompt.yaml. Do not edit manually. -->

## Purpose

Rules for automatic task-mode selection by task intent.

## Persona

- The agent classifies task type before acting.
- The agent chooses the minimal sufficient mode.
- The agent records mode choice in spec_changes.

## Workflow Steps

### 1. Task Classification

1. bugfix: fix a defect without expanding scope/functionality.
2. development: implement new functionality or behavior changes.
3. refactor: improve structure/quality without intended behavior changes.
4. tests: add or fix tests as the primary goal.
5. analysis: answer questions about the system without file changes.

### 2. Mode Application

1. Select exactly one primary mode.
2. If secondary mode is needed, record it as secondary in spec_changes.
3. Apply workflow from .spec_agent/prompts/modes/.

## Must Rules

1. Default to bugfix for defect-fix requests unless redesign is requested.
2. development/refactor/tests require iterative plan-based execution.
3. analysis mode must not change code or specifications.
4. Every mode requires creating/updating spec_changes file.

## Tool Policy

### Allowed

- Use request context and local evidence to classify task type.

### Forbidden

- Choosing a heavier mode without reason.
- Skipping mode record in spec_changes.

### Notes

- If type is ambiguous, choose safest mode and state assumptions.

## Output Contract

### Required Sections

- Task type
- Selected mode
- Why this mode
- spec_changes file

### Rules

- Mode choice must be concise and verifiable.

