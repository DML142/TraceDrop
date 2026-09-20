# TraceDrop

TraceDrop is a lightweight developer observability tool for collecting application traces, inspecting event timelines, simulating failures and retries, and viewing basic reliability metrics.

## Development status

The project is currently in active MVP development.

Implemented:
- Go API bootstrap and health endpoint;
- Next.js dashboard shell;
- Docker development environment;
- PostgreSQL schema and migrations;
- pgx connection pool;
- sqlc query definitions;
- trace persistence repository;
- PostgreSQL integration tests;
- GitHub Actions CI.

## Local database

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
