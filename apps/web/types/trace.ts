export type TraceStatus = "pending" | "running" | "completed" | "failed";

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

export type TraceListResponse = {
  traces: Trace[];
};
