# TraceDrop — Agent Engineering Rules

1. Read tech.md, AGENTS.md and README.md before changing code.
2. Implement only the current phase or explicitly requested task.
3. Prefer simple architecture and small dependencies.
4. Production-minded MVP: validation, errors, context, logging and tests still matter.
5. Go: prefer stdlib, context.Context, errors.Is/errors.As, slog, server timeouts and graceful shutdown.
6. Never silently ignore meaningful errors.
7. Keep handlers thin; do not mix transport, business logic and persistence.
8. Comments should explain non-obvious reasoning, not restate code.
9. Prefer intent-revealing names; avoid vague Helper/Utils/Manager names.
10. PostgreSQL is authoritative; all schema changes use migrations.
11. Use transactions only when atomicity is required.
12. API errors must be consistent and must not expose internals.
13. Frontend must handle loading, empty, error and success states.
14. Avoid any and unnecessary global state in TypeScript.
15. Test behavior, not implementation details.
16. Do not weaken CI to make a PR pass.
17. Never commit secrets or log sensitive values.
18. Add dependencies only when they solve a concrete need.
19. Work on feature branches, never directly on main.
20. Keep commits coherent and PRs reviewable.
21. Before committing, run relevant format/lint/typecheck/test/build checks.
22. Explain material architecture changes before making them.
23. Keep docs synchronized with implementation.
24. Never claim completion without verification.

PR format:

## What changed
- ...

## Testing
- ...

Do not mention AI tooling in branch names, commits or PR descriptions unless explicitly requested.
