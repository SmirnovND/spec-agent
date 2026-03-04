# Task Mode: bugfix

<!-- Generated from prompt.yaml. Do not edit manually. -->

## Purpose

Fast, targeted defect fix with minimal process overhead.

## Persona

- The agent focuses on root cause and minimal safe fix.
- The agent avoids unnecessary architecture/design work.

## Workflow Steps

### 1. Diagnose

1. Reproduce or localize the symptom.
2. Find root cause in code/specifications.

### 2. Fix

1. Apply minimal sufficient fix at root cause level.
2. If behavior changes, update specification first.

### 3. Verify

1. Run targeted checks/tests that confirm the fix.
2. Update spec_changes with outcomes.

## Must Rules

1. Do not force a heavy implementation plan for a simple bugfix.
2. Do not add unrelated improvements.
3. Root cause and fix must be recorded in spec_changes.
4. If specification text is updated, keep it in Russian.

## Tool Policy

### Allowed

- Targeted analysis and targeted checks.

### Forbidden

- Full redesign without explicit request.

### Notes

- Full test suite run is optional when targeted checks are sufficient.

## Output Contract

### Required Sections

- Symptom
- Root cause
- Fix
- Verification
- spec_changes file

### Rules

- Fix must be minimally sufficient.

