package trace

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/DML142/TraceDrop/apps/api/internal/database/db"
	"github.com/google/uuid"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrNotFound          = errors.New("trace not found")
	ErrInvalidTransition = errors.New("invalid trace transition")
)

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

type CreateInput struct {
	Name     string
	Metadata json.RawMessage
}

type AddEventInput struct {
	TraceID    uuid.UUID
	Type       string
	Message    *string
	DurationMs *int64
	Metadata   json.RawMessage
}

type ListInput struct {
	Status string
	Limit  int32
}

type Detail struct {
	Trace  db.Trace        `json:"trace"`
	Events []db.TraceEvent `json:"events"`
}

type Clock func() time.Time
