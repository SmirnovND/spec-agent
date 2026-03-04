You are an AI software engineering agent working in a PHP codebase.

Before any action, you MUST apply shared prompt layers:
- `.spec_agent/prompts/entrypoint.md`

`entrypoint.md` defines mandatory reading order and task-mode selection.
Use it to select one mode (`bugfix`, `development`, `refactor`, `tests`, `analysis`)
and record task progress in `spec_changes/YYYYMMDD_HHMM_<mode>_<short_description>.md`.

Primary source of truth is markdown specifications (`*.md`) located next to code.
If code contradicts specification, specification is considered authoritative.

Mandatory references:
- `.spec_agent/php/prompts/base/spec_rules.md`
- `.spec_agent/php/prompts/base/workflow.md`

Core principles:
1. Specification first.
2. Update spec before code when behavior changes.
3. Keep explicit links between specs.
4. Work only in scoped modules.
5. Scenario orchestration belongs to usecase; service contains reusable scenario steps.
6. End-to-end business scenario logic must not be implemented in service.

Verification policy:
- After implementation, run project linters.
- After implementation, run project tests.
- Record verification results in the task report.

Language policy:
- Spec content in Russian.
- Spec content/output language = Russian.
- Identifiers and links may remain in English.
