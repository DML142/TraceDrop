-- name: CreateTrace :one
INSERT INTO traces (name, status, metadata)
VALUES ($1, $2, $3)
RETURNING id, name, status, started_at, finished_at, duration_ms, metadata, created_at;

-- name: GetTrace :one
SELECT id, name, status, started_at, finished_at, duration_ms, metadata, created_at
FROM traces
WHERE id = $1;

-- name: ListTraces :many
SELECT id, name, status, started_at, finished_at, duration_ms, metadata, created_at
FROM traces
ORDER BY created_at DESC
LIMIT $1;

-- name: ListTracesByStatus :many
SELECT id, name, status, started_at, finished_at, duration_ms, metadata, created_at
FROM traces
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: TransitionTrace :one
UPDATE traces
SET status = $3,
    finished_at = $4,
    duration_ms = $5
WHERE id = $1
  AND status = $2
RETURNING id, name, status, started_at, finished_at, duration_ms, metadata, created_at;

-- name: CreateTraceEvent :one
INSERT INTO trace_events (trace_id, event_type, message, duration_ms, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, trace_id, event_type, message, timestamp, duration_ms, metadata, created_at;

-- name: ListTraceEvents :many
SELECT id, trace_id, event_type, message, timestamp, duration_ms, metadata, created_at
FROM trace_events
WHERE trace_id = $1
ORDER BY timestamp ASC, created_at ASC;
