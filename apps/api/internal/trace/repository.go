package trace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/google/uuid"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) Create(ctx context.Context, name string, metadata json.RawMessage) (db.Trace, error) {
	item, err := r.queries.CreateTrace(ctx, db.CreateTraceParams{
		Name:     name,
		Status:   "running",
		Metadata: metadata,
	})
	if err != nil {
		return db.Trace{}, fmt.Errorf("create trace: %w", err)
	}
	return item, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (db.Trace, error) {
	item, err := r.queries.GetTrace(ctx, id)
	if err != nil {
		return db.Trace{}, fmt.Errorf("get trace: %w", err)
	}
	return item, nil
}

func (r *Repository) List(ctx context.Context, limit int32) ([]db.Trace, error) {
	items, err := r.queries.ListTraces(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list traces: %w", err)
	}
	return items, nil
}

func (r *Repository) AddEvent(
	ctx context.Context,
	traceID uuid.UUID,
	eventType string,
	message *string,
	durationMs *int64,
	metadata json.RawMessage,
) (db.TraceEvent, error) {
	item, err := r.queries.CreateTraceEvent(ctx, db.CreateTraceEventParams{
		TraceID:    traceID,
		EventType:  eventType,
		Message:    message,
		DurationMs: durationMs,
		Metadata:   metadata,
	})
	if err != nil {
		return db.TraceEvent{}, fmt.Errorf("create trace event: %w", err)
	}
	return item, nil
}

func (r *Repository) Events(ctx context.Context, traceID uuid.UUID) ([]db.TraceEvent, error) {
	items, err := r.queries.ListTraceEvents(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("list trace events: %w", err)
	}
	return items, nil
}
