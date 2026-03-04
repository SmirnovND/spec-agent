# Task Mode: refactor

<!-- Generated from prompt.yaml. Do not edit manually. -->

## Purpose

Improve code structure/quality without intended behavior change.

## Persona

- The agent performs pre-refactor analysis first.
- The agent executes refactoring in safe, iterative steps.

## Workflow Steps

### 1. Pre-Refactor Analysis

1. Record refactor drivers (duplication, complexity, coupling, debt).
2. Define constraints and invariants that must not change.

### 2. Refactor Plan

1. Split refactor into safe incremental steps.
2. Define expected behavior invariants.

### 3. Iterative Refactoring

1. Execute one refactor step at a time.
2. Validate invariants after each step.

### 4. Verification

1. Run relevant tests/linters.
2. Record behavior-preservation evidence in spec_changes.

## Must Rules

1. Refactor requires analysis and a plan.
2. Do not change business behavior unless switching to development mode.
3. Behavior invariants must be explicit in spec_changes.
4. If specification text is updated, keep it in Russian.

## Tool Policy

### Allowed

- Structural improvements that increase maintainability.

### Forbidden

- Hidden functional changes under refactor label.

### Notes

- If behavior change is required, explicitly switch mode and record it.

## Output Contract

### Required Sections

- Refactor goals
- Analysis
- Plan
- Invariants
- Verification
- spec_changes file

### Rules

- Final report must confirm no unintended behavior change.

