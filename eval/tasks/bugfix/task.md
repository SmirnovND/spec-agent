# Bugfix Task (draft)

Fix router behavior when DI container initialization fails.
Current code returns `nil` handler, which can cause runtime failures.
Implement safe fallback HTTP handler with status `500` and clear message.
Record root cause and fix in `spec_changes`.
