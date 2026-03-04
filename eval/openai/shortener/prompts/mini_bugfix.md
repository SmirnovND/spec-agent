You are working on a Go backend in this workspace.

Before coding:
- Read `README.md`.
- Read `.spec_agent/prompts/entrypoint.md` and follow the prompt stack described there.
- Read baseline spec `internal/router/router.md`.

Task (small but strict):
Fix router behavior on DI failure.

Required changes:
1. In `internal/router/router.go`, when DI initialization fails, do NOT return `nil`.
2. Return a fallback `http.HandlerFunc` that responds with `500` and message `Service initialization failed`.
3. Add a focused test in `internal/router/router_fallback_test.go`:
   - test name: `TestHandlerFallbackOnDIError`
   - verify returned handler is non-nil and responds with status 500.
4. Update `internal/router/router.md` to document fallback handler behavior.
5. Add changelog note file in `spec_changes/` with `bugfix` in filename.

Scope constraint:
- Change only:
  - `internal/router/**`
  - `spec_changes/**`
  - `.spec_agent/**`
  - `README.md`

After completion write `.eval/openai/mini_bugfix_report.json`:
{
  "task": "mini_bugfix_router_fallback",
  "status": "done",
  "changed_files": ["..."],
  "notes": "..."
}
