package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/DML142/TraceDrop/apps/api/internal/trace"
	"github.com/google/uuid"
)

type TraceService interface {
	Create(context.Context, trace.CreateInput) (db.Trace, error)
	Get(context.Context, uuid.UUID) (trace.Detail, error)
	List(context.Context, trace.ListInput) ([]db.Trace, error)
	AddEvent(context.Context, trace.AddEventInput) (db.TraceEvent, error)
	Complete(context.Context, uuid.UUID) (db.Trace, error)
	Fail(context.Context, uuid.UUID) (db.Trace, error)
}

type API struct {
	traces TraceService
	ready  func(context.Context) error
}

func New(traces TraceService, ready func(context.Context) error) *API {
	return &API{traces: traces, ready: ready}
}

func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.health)
	mux.HandleFunc("GET /ready", a.readiness)
	mux.HandleFunc("POST /api/v1/traces", a.createTrace)
	mux.HandleFunc("GET /api/v1/traces", a.listTraces)
	mux.HandleFunc("GET /api/v1/traces/{traceId}", a.getTrace)
	mux.HandleFunc("POST /api/v1/traces/{traceId}/events", a.addEvent)
	mux.HandleFunc("POST /api/v1/traces/{traceId}/complete", a.completeTrace)
	mux.HandleFunc("POST /api/v1/traces/{traceId}/fail", a.failTrace)
	return mux
}

func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) readiness(w http.ResponseWriter, r *http.Request) {
	if a.ready != nil {
		if err := a.ready(r.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "not_ready", "Service is not ready")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a *API) createTrace(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string          `json:"name"`
		Metadata json.RawMessage `json:"metadata"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Request body must be valid JSON")
		return
	}

	item, err := a.traces.Create(r.Context(), trace.CreateInput{Name: body.Name, Metadata: body.Metadata})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) listTraces(w http.ResponseWriter, r *http.Request) {
	var limit int32
	if raw := r.URL.Query().Get("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "limit must be an integer")
			return
		}
		limit = int32(value)
	}

	items, err := a.traces.List(r.Context(), trace.ListInput{
		Status: r.URL.Query().Get("status"),
		Limit:  limit,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"traces": items})
}

func (a *API) getTrace(w http.ResponseWriter, r *http.Request) {
	id, ok := traceID(w, r)
	if !ok {
		return
	}

	item, err := a.traces.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) addEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := traceID(w, r)
	if !ok {
		return
	}

	var body struct {
		Type       string          `json:"type"`
		Message    *string         `json:"message"`
		DurationMs *int64          `json:"duration_ms"`
		Metadata   json.RawMessage `json:"metadata"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Request body must be valid JSON")
		return
	}

	item, err := a.traces.AddEvent(r.Context(), trace.AddEventInput{
		TraceID:    id,
		Type:       body.Type,
		Message:    body.Message,
		DurationMs: body.DurationMs,
		Metadata:   body.Metadata,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) completeTrace(w http.ResponseWriter, r *http.Request) {
	a.transition(w, r, a.traces.Complete)
}

func (a *API) failTrace(w http.ResponseWriter, r *http.Request) {
	a.transition(w, r, a.traces.Fail)
}

func (a *API) transition(
	w http.ResponseWriter,
	r *http.Request,
	action func(context.Context, uuid.UUID) (db.Trace, error),
) {
	id, ok := traceID(w, r)
	if !ok {
		return
	}

	item, err := action(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func traceID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("traceId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "traceId must be a valid UUID")
		return uuid.Nil, false
	}
	return id, true
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return errors.New("request body must contain one JSON value")
	}
	return nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, trace.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, trace.ErrNotFound):
		writeError(w, http.StatusNotFound, "trace_not_found", "Trace was not found")
	case errors.Is(err, trace.ErrInvalidTransition):
		writeError(w, http.StatusConflict, "invalid_transition", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
