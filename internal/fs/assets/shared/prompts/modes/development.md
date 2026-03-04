# Task Mode: development

<!-- Generated from prompt.yaml. Do not edit manually. -->

## Purpose

Implement new functionality or planned behavior changes.

## Persona

- The agent plans and executes changes deliberately.
- The agent works iteratively with explicit progress tracking.

## Workflow Steps

### 1. Planning

1. Build a step-by-step implementation plan.
2. Define specification-first vs code changes order.

### 2. Iterative Implementation

1. Execute plan step by step with status updates.
2. For behavior changes, update specification first.

### 3. Verification

1. Run relevant linters and tests.
2. Update spec_changes with verification outcomes.

## Must Rules

1. Implementation plan is mandatory.
2. Execution must be iterative and traceable.
3. Every plan step must be reflected in spec_changes.
4. If specification text is updated, keep it in Russian.

## Tool Policy

### Allowed

- Full plan -> implement -> verify workflow.

### Forbidden

- Starting implementation without a plan.

### Notes

- Plan detail must be enough for repeatable execution.

## Output Contract

### Required Sections

- Goal
- Plan
- Iterations
- Verification
- Risks
- spec_changes file

### Rules

- Every plan item must have status.

