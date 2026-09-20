package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/google/uuid"
)

type Store interface {
	Create(context.Context, string, json.RawMessage) (db.Trace, error)
	Get(context.Context, uuid.UUID) (db.Trace, error)
	List(context.Context, int32) ([]db.Trace, error)
	ListByStatus(context.Context, string, int32) ([]db.Trace, error)
	AddEvent(context.Context, uuid.UUID, string, *string, *int64, json.RawMessage) (db.TraceEvent, error)
	Events(context.Context, uuid.UUID) ([]db.TraceEvent, error)
	Transition(context.Context, uuid.UUID, string, string, *time.Time, *int64) (db.Trace, error)
}

type Service struct {
	store Store
	now   Clock
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
		now:   time.Now,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (db.Trace, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > 120 {
		return db.Trace{}, fmt.Errorf("%w: name must be between 1 and 120 characters", ErrInvalidInput)
	}

	metadata, err := normalizeMetadata(input.Metadata)
	if err != nil {
		return db.Trace{}, err
	}

	return s.store.Create(ctx, name, metadata)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (Detail, error) {
	item, err := s.store.Get(ctx, id)
	if err != nil {
		return Detail{}, err
	}

	events, err := s.store.Events(ctx, id)
	if err != nil {
		return Detail{}, err
	}

	return Detail{Trace: item, Events: events}, nil
}

func (s *Service) List(ctx context.Context, input ListInput) ([]db.Trace, error) {
	limit := input.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("%w: limit must be between 1 and 100", ErrInvalidInput)
	}

	if input.Status == "" {
		return s.store.List(ctx, limit)
	}
	if !validStatus(input.Status) {
		return nil, fmt.Errorf("%w: unsupported trace status", ErrInvalidInput)
	}

	return s.store.ListByStatus(ctx, input.Status, limit)
}

func (s *Service) AddEvent(ctx context.Context, input AddEventInput) (db.TraceEvent, error) {
	if !validEventType(input.Type) {
		return db.TraceEvent{}, fmt.Errorf("%w: unsupported event type", ErrInvalidInput)
	}
	if input.DurationMs != nil && *input.DurationMs < 0 {
		return db.TraceEvent{}, fmt.Errorf("%w: duration_ms cannot be negative", ErrInvalidInput)
	}

	metadata, err := normalizeMetadata(input.Metadata)
	if err != nil {
		return db.TraceEvent{}, err
	}

	item, err := s.store.Get(ctx, input.TraceID)
	if err != nil {
		return db.TraceEvent{}, err
	}
	if item.Status == StatusCompleted {
		return db.TraceEvent{}, fmt.Errorf("%w: completed traces cannot receive events", ErrInvalidTransition)
	}

	return s.store.AddEvent(ctx, input.TraceID, input.Type, input.Message, input.DurationMs, metadata)
}

func (s *Service) Complete(ctx context.Context, id uuid.UUID) (db.Trace, error) {
	return s.finish(ctx, id, StatusCompleted)
}

func (s *Service) Fail(ctx context.Context, id uuid.UUID) (db.Trace, error) {
	return s.finish(ctx, id, StatusFailed)
}

func (s *Service) finish(ctx context.Context, id uuid.UUID, nextStatus string) (db.Trace, error) {
	item, err := s.store.Get(ctx, id)
	if err != nil {
		return db.Trace{}, err
	}
	if item.Status != StatusRunning {
		return db.Trace{}, fmt.Errorf("%w: %s cannot transition to %s", ErrInvalidTransition, item.Status, nextStatus)
	}

	now := s.now().UTC()
	duration := now.Sub(item.StartedAt).Milliseconds()
	if duration < 0 {
		duration = 0
	}

	return s.store.Transition(ctx, id, StatusRunning, nextStatus, &now, &duration)
}

func normalizeMetadata(metadata json.RawMessage) (json.RawMessage, error) {
	if len(metadata) == 0 {
		return json.RawMessage(`{}`), nil
	}

	var value map[string]any
	if err := json.Unmarshal(metadata, &value); err != nil {
		return nil, fmt.Errorf("%w: metadata must be a JSON object", ErrInvalidInput)
	}
	return metadata, nil
}

func validStatus(status string) bool {
	switch status {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

func validEventType(eventType string) bool {
	switch eventType {
	case "received", "processing", "retry", "completed", "failed", "custom":
		return true
	default:
		return false
	}
}
