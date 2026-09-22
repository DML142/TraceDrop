export type TraceStatus = "pending" | "running" | "completed" | "failed";

export type TraceEventType =
  | "received"
  | "processing"
  | "retry"
  | "completed"
  | "failed"
  | "custom";

export type Trace = {
  id: string;
  name: string;
  status: TraceStatus;
  started_at: string;
  finished_at: string | null;
  duration_ms: number | null;
  metadata: Record<string, unknown>;
  created_at: string;
};

export type TraceEvent = {
  id: string;
  trace_id: string;
  event_type: TraceEventType;
  message: string | null;
  timestamp: string;
  duration_ms: number | null;
  metadata: Record<string, unknown>;
  created_at: string;
};

export type TraceListResponse = {
  traces: Trace[];
};

export type TraceDetailResponse = {
  trace: Trace;
  events: TraceEvent[];
};
