-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE traces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ,
    duration_ms BIGINT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_traces_created_at ON traces (created_at DESC);
CREATE INDEX idx_traces_status ON traces (status);

CREATE TABLE trace_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trace_id UUID NOT NULL REFERENCES traces(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL CHECK (event_type IN ('received', 'processing', 'retry', 'completed', 'failed', 'custom')),
    message TEXT,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
    duration_ms BIGINT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_trace_events_trace_timestamp
    ON trace_events (trace_id, timestamp);

-- +goose Down
DROP TABLE IF EXISTS trace_events;
DROP TABLE IF EXISTS traces;
