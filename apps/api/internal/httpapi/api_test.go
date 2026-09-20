package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/DML142/TraceDrop/apps/api/internal/trace"
	"github.com/google/uuid"
)

type stubService struct {
	createFn   func(context.Context, trace.CreateInput) (db.Trace, error)
	getFn      func(context.Context, uuid.UUID) (trace.Detail, error)
	listFn     func(context.Context, trace.ListInput) ([]db.Trace, error)
	addEventFn func(context.Context, trace.AddEventInput) (db.TraceEvent, error)
	completeFn func(context.Context, uuid.UUID) (db.Trace, error)
	failFn     func(context.Context, uuid.UUID) (db.Trace, error)
}

func (s stubService) Create(ctx context.Context, input trace.CreateInput) (db.Trace, error) {
	return s.createFn(ctx, input)
}
func (s stubService) Get(ctx context.Context, id uuid.UUID) (trace.Detail, error) {
	return s.getFn(ctx, id)
}
func (s stubService) List(ctx context.Context, input trace.ListInput) ([]db.Trace, error) {
	return s.listFn(ctx, input)
}
func (s stubService) AddEvent(ctx context.Context, input trace.AddEventInput) (db.TraceEvent, error) {
	return s.addEventFn(ctx, input)
}
func (s stubService) Complete(ctx context.Context, id uuid.UUID) (db.Trace, error) {
	return s.completeFn(ctx, id)
}
func (s stubService) Fail(ctx context.Context, id uuid.UUID) (db.Trace, error) {
	return s.failFn(ctx, id)
}

func defaultStub() stubService {
	return stubService{
		createFn:   func(context.Context, trace.CreateInput) (db.Trace, error) { return db.Trace{}, nil },
		getFn:      func(context.Context, uuid.UUID) (trace.Detail, error) { return trace.Detail{}, nil },
		listFn:     func(context.Context, trace.ListInput) ([]db.Trace, error) { return nil, nil },
		addEventFn: func(context.Context, trace.AddEventInput) (db.TraceEvent, error) { return db.TraceEvent{}, nil },
		completeFn: func(context.Context, uuid.UUID) (db.Trace, error) { return db.Trace{}, nil },
		failFn:     func(context.Context, uuid.UUID) (db.Trace, error) { return db.Trace{}, nil },
	}
}

func TestCreateTraceReturnsCreated(t *testing.T) {
	id := uuid.New()
	service := defaultStub()
	service.createFn = func(_ context.Context, input trace.CreateInput) (db.Trace, error) {
		return db.Trace{
			ID:        id,
			Name:      input.Name,
			Status:    trace.StatusRunning,
			StartedAt: time.Now().UTC(),
			Metadata:  json.RawMessage(`{}`),
		}, nil
	}

	api := New(service, nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/traces", bytes.NewBufferString(`{"name":"payment"}`))
	recorder := httptest.NewRecorder()

	api.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestInvalidTraceIDReturnsBadRequest(t *testing.T) {
	api := New(defaultStub(), nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/traces/not-a-uuid", nil)
	recorder := httptest.NewRecorder()

	api.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestMissingTraceReturnsNotFound(t *testing.T) {
	service := defaultStub()
	service.getFn = func(context.Context, uuid.UUID) (trace.Detail, error) {
		return trace.Detail{}, trace.ErrNotFound
	}

	api := New(service, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/traces/"+uuid.NewString(), nil)
	recorder := httptest.NewRecorder()

	api.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recorder.Code)
	}
}

func TestInvalidTransitionReturnsConflict(t *testing.T) {
	service := defaultStub()
	service.completeFn = func(context.Context, uuid.UUID) (db.Trace, error) {
		return db.Trace{}, trace.ErrInvalidTransition
	}

	api := New(service, nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/traces/"+uuid.NewString()+"/complete", nil)
	recorder := httptest.NewRecorder()

	api.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", recorder.Code)
	}
}

func TestReadinessFailureReturnsServiceUnavailable(t *testing.T) {
	api := New(defaultStub(), func(context.Context) error {
		return errors.New("database unavailable")
	})
	request := httptest.NewRequest(http.MethodGet, "/ready", nil)
	recorder := httptest.NewRecorder()

	api.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", recorder.Code)
	}
}
