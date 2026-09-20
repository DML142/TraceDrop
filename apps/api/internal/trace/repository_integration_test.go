package trace

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/DML142/TraceDrop/apps/api/internal/database"
	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
)

func TestRepositoryTracePersistence(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := database.Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, "TRUNCATE trace_events, traces RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	repository := NewRepository(db.New(pool))

	created, err := repository.Create(ctx, "payment-processing", json.RawMessage(`{"order_id":"ord_123"}`))
	if err != nil {
		t.Fatalf("create trace: %v", err)
	}

	got, err := repository.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get trace: %v", err)
	}
	if got.Name != created.Name {
		t.Fatalf("expected trace name %q, got %q", created.Name, got.Name)
	}
	if got.Status != "running" {
		t.Fatalf("expected running status, got %q", got.Status)
	}

	message := "received payment request"
	event, err := repository.AddEvent(
		ctx,
		created.ID,
		"received",
		&message,
		nil,
		json.RawMessage(`{"source":"integration-test"}`),
	)
	if err != nil {
		t.Fatalf("add event: %v", err)
	}
	if event.TraceID != created.ID {
		t.Fatalf("event trace id does not match")
	}

	events, err := repository.Events(ctx, created.ID)
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	traces, err := repository.List(ctx, 10)
	if err != nil {
		t.Fatalf("list traces: %v", err)
	}
	if len(traces) != 1 {
		t.Fatalf("expected 1 trace, got %d", len(traces))
	}
}
