package trace

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/google/uuid"
)

type fakeStore struct {
	trace      db.Trace
	event      db.TraceEvent
	traces     []db.Trace
	events     []db.TraceEvent
	getErr     error
	transition struct {
		called   bool
		from     string
		to       string
		finished *time.Time
		duration *int64
	}
}

func (f *fakeStore) Create(_ context.Context, name string, metadata json.RawMessage) (db.Trace, error) {
	item := f.trace
	item.Name = name
	item.Metadata = metadata
	if item.Status == "" {
		item.Status = StatusRunning
	}
	return item, nil
}

func (f *fakeStore) Get(_ context.Context, _ uuid.UUID) (db.Trace, error) {
	if f.getErr != nil {
		return db.Trace{}, f.getErr
	}
	return f.trace, nil
}

func (f *fakeStore) List(_ context.Context, _ int32) ([]db.Trace, error) {
	return f.traces, nil
}

func (f *fakeStore) ListByStatus(_ context.Context, _ string, _ int32) ([]db.Trace, error) {
	return f.traces, nil
}

func (f *fakeStore) AddEvent(_ context.Context, _ uuid.UUID, _ string, _ *string, _ *int64, _ json.RawMessage) (db.TraceEvent, error) {
	return f.event, nil
}

func (f *fakeStore) Events(_ context.Context, _ uuid.UUID) ([]db.TraceEvent, error) {
	return f.events, nil
}

func (f *fakeStore) Transition(_ context.Context, _ uuid.UUID, from, to string, finished *time.Time, duration *int64) (db.Trace, error) {
	f.transition.called = true
	f.transition.from = from
	f.transition.to = to
	f.transition.finished = finished
	f.transition.duration = duration

	item := f.trace
	item.Status = to
	item.FinishedAt = finished
	item.DurationMs = duration
	return item, nil
}

func TestCreateRejectsEmptyName(t *testing.T) {
	service := NewService(&fakeStore{})

	_, err := service.Create(context.Background(), CreateInput{Name: "   "})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestListRejectsUnsupportedStatus(t *testing.T) {
	service := NewService(&fakeStore{})

	_, err := service.List(context.Background(), ListInput{Status: "unknown"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAddEventRejectsCompletedTrace(t *testing.T) {
	store := &fakeStore{trace: db.Trace{
		ID:     uuid.New(),
		Status: StatusCompleted,
	}}
	service := NewService(store)

	_, err := service.AddEvent(context.Background(), AddEventInput{
		TraceID: store.trace.ID,
		Type:    "custom",
	})
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got %v", err)
	}
}

func TestCompleteCalculatesDurationAndTransitions(t *testing.T) {
	startedAt := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
	now := startedAt.Add(1500 * time.Millisecond)
	store := &fakeStore{trace: db.Trace{
		ID:        uuid.New(),
		Status:    StatusRunning,
		StartedAt: startedAt,
	}}
	service := NewService(store)
	service.now = func() time.Time { return now }

	item, err := service.Complete(context.Background(), store.trace.ID)
	if err != nil {
		t.Fatalf("complete trace: %v", err)
	}
	if !store.transition.called {
		t.Fatal("expected transition to be called")
	}
	if store.transition.from != StatusRunning || store.transition.to != StatusCompleted {
		t.Fatalf("unexpected transition %s -> %s", store.transition.from, store.transition.to)
	}
	if store.transition.duration == nil || *store.transition.duration != 1500 {
		t.Fatalf("expected 1500ms duration, got %v", store.transition.duration)
	}
	if item.Status != StatusCompleted {
		t.Fatalf("expected completed status, got %q", item.Status)
	}
}

func TestCompleteRejectsNonRunningTrace(t *testing.T) {
	store := &fakeStore{trace: db.Trace{
		ID:        uuid.New(),
		Status:    StatusFailed,
		StartedAt: time.Now(),
	}}
	service := NewService(store)

	_, err := service.Complete(context.Background(), store.trace.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got %v", err)
	}
}
