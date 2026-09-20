# TraceDrop — Technical Plan

## Product
TraceDrop is a lightweight developer observability tool for inspecting application events and execution traces.

MVP capabilities:
- create traces;
- append trace events;
- list/filter traces;
- inspect a trace timeline;
- calculate basic latency/error metrics;
- simulate success/failure/retry scenarios;
- health/readiness endpoints;
- tests, Docker, CI/CD, preview deployments.

## Stack
Backend: Go, net/http, PostgreSQL, pgx, sqlc, goose, slog.
Frontend: Next.js, TypeScript, App Router, Tailwind CSS, TanStack Query.
Infra: Docker, Docker Compose, GitHub Actions, Vercel for web, managed API host and PostgreSQL.

## Backend architecture
HTTP -> Handler -> Service -> Repository -> PostgreSQL

Handlers own transport concerns.
Services own business rules and state transitions.
Repositories own persistence.

## Trace statuses
pending
running
completed
failed

Current MVP transitions:
running -> completed
running -> failed

Retry transitions will be added with the simulator workflow.

## API surface
GET /health
GET /ready
POST /api/v1/traces
GET /api/v1/traces
GET /api/v1/traces/{traceId}
POST /api/v1/traces/{traceId}/events
POST /api/v1/traces/{traceId}/complete
POST /api/v1/traces/{traceId}/fail
GET /api/v1/metrics

## Development phases
- [x] Phase 0: bootstrap Go API, Next.js web, Docker Compose, env example, README, CI skeleton, health endpoint.
- [x] Phase 1: PostgreSQL, migrations, pgx/sqlc, repositories, integration tests.
- [x] Phase 2: trace API and domain rules.
- [ ] Phase 3: dashboard and filtering. Current PR.
- [ ] Phase 4: trace timeline.
- [ ] Phase 5: simulator.
- [ ] Phase 6: metrics.
- [ ] Phase 7: production hardening.
- [ ] Phase 8: full automated test coverage.
- [ ] Phase 9: deployment and previews.
- [ ] Phase 10: portfolio polish and architecture docs.

Do not expand scope before the MVP phases are complete.
