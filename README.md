# TraceDrop

TraceDrop is a lightweight developer observability tool for collecting application traces, inspecting event timelines, simulating failures and retries, and viewing basic reliability metrics.

## Development status

Implemented:
- Go API and health/readiness endpoints;
- Next.js dashboard shell;
- Docker development environment;
- PostgreSQL schema and migrations;
- pgx/sqlc persistence;
- trace repository and integration tests;
- trace HTTP API with validation and lifecycle rules;
- GitHub Actions CI.

## Trace API

```text
POST /api/v1/traces
GET  /api/v1/traces
GET  /api/v1/traces/{traceId}
POST /api/v1/traces/{traceId}/events
POST /api/v1/traces/{traceId}/complete
POST /api/v1/traces/{traceId}/fail
```

List traces accepts `status` and `limit` query parameters.

## Local development

Start PostgreSQL:

```bash
docker compose up -d postgres
```

Apply migrations:

```bash
cd apps/api
make migrate-up
```

Run backend tests:

```bash
make test
```
