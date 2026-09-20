"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";

import { getTrace } from "@/lib/api";

import { MetadataBlock } from "./metadata-block";
import { TraceTimeline } from "./trace-timeline";

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "medium",
  }).format(new Date(value));
}

function formatDuration(duration: number | null) {
  if (duration === null) {
    return "Running";
  }
  if (duration < 1000) {
    return `${duration} ms`;
  }
  return `${(duration / 1000).toFixed(2)} s`;
}

export function TraceDetail({ traceId }: { traceId: string }) {
  const detail = useQuery({
    queryKey: ["trace", traceId],
    queryFn: () => getTrace(traceId),
  });

  if (detail.isPending) {
    return (
      <main className="shell">
        <div className="state">Loading trace…</div>
      </main>
    );
  }

  if (detail.isError) {
    return (
      <main className="shell">
        <Link className="back-link" href="/">
          ← Back to traces
        </Link>
        <div className="panel state state-error">
          <strong>Could not load trace</strong>
          <span>The trace may not exist or the API may be unavailable.</span>
        </div>
      </main>
    );
  }

  const { trace, events } = detail.data;

  return (
    <main className="shell">
      <Link className="back-link" href="/">
        ← Back to traces
      </Link>

      <header className="trace-detail-header">
        <div>
          <div className="eyebrow">Trace detail</div>
          <h1 className="trace-detail-title">{trace.name}</h1>
          <div className="trace-id trace-detail-id">{trace.id}</div>
        </div>
        <span className={`status status-${trace.status}`}>{trace.status}</span>
      </header>

      <section className="trace-summary">
        <article>
          <span>Started</span>
          <strong>{formatDate(trace.started_at)}</strong>
        </article>
        <article>
          <span>Duration</span>
          <strong>{formatDuration(trace.duration_ms)}</strong>
        </article>
        <article>
          <span>Events</span>
          <strong>{events.length}</strong>
        </article>
      </section>

      <div className="trace-detail-grid">
        <section className="panel timeline-panel">
          <div className="panel-header">
            <div>
              <h2>Timeline</h2>
              <p>Ordered events recorded for this trace.</p>
            </div>
          </div>
          <TraceTimeline events={events} />
        </section>

        <aside className="detail-sidebar">
          <MetadataBlock title="Trace metadata" metadata={trace.metadata} />
        </aside>
      </div>
    </main>
  );
}
