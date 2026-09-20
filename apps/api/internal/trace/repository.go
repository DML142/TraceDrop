package trace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) Create(ctx context.Context, name string, metadata json.RawMessage) (db.Trace, error) {
	item, err := r.queries.CreateTrace(ctx, db.CreateTraceParams{Name: name, Status: StatusRunning, Metadata: metadata})
	if err != nil {
		return db.Trace{}, fmt.Errorf("create trace: %w", err)
	}
	return item, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (db.Trace, error) {
	item, err := r.queries.GetTrace(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Trace{}, ErrNotFound
	}
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

func (r *Repository) ListByStatus(ctx context.Context, status string, limit int32) ([]db.Trace, error) {
	items, err := r.queries.ListTracesByStatus(ctx, db.ListTracesByStatusParams{Status: status, Limit: limit})
	if err != nil {
		return nil, fmt.Errorf("list traces by status: %w", err)
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
		TraceID: traceID, EventType: eventType, Message: message, DurationMs: durationMs, Metadata: metadata,
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

func (r *Repository) Transition(
	ctx context.Context,
	id uuid.UUID,
	from string,
	to string,
	finishedAt *time.Time,
	durationMs *int64,
) (db.Trace, error) {
	item, err := r.queries.TransitionTrace(ctx, db.TransitionTraceParams{
		ID: id, Status: from, Status_2: to, FinishedAt: finishedAt, DurationMs: durationMs,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return db.Trace{}, ErrInvalidTransition
	}
	if err != nil {
		return db.Trace{}, fmt.Errorf("transition trace: %w", err)
	}
	return item, nil
}
