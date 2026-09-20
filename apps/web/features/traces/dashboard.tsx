"use client";

import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

import { getTraces } from "@/lib/api";
import type { TraceStatus } from "@/types/trace";

import { StatusFilter } from "./status-filter";
import { TraceTable } from "./trace-table";

export function Dashboard() {
  const [status, setStatus] = useState<TraceStatus | undefined>();

  const traces = useQuery({
    queryKey: ["traces", status ?? "all"],
    queryFn: () => getTraces(status),
  });

  return (
    <main className="shell">
      <header className="page-header">
        <div>
          <div className="eyebrow">Developer observability</div>
          <h1>TraceDrop</h1>
          <p className="subtitle">Inspect application traces and execution outcomes.</p>
        </div>
        <div className="live-indicator">
          <span className="live-dot" />
          API connected
        </div>
      </header>

      <section className="panel">
        <div className="panel-header">
          <div>
            <h2>Recent traces</h2>
            <p>Latest trace activity from the ingestion API.</p>
          </div>
          <StatusFilter value={status} onChange={setStatus} />
        </div>

        {traces.isPending ? <div className="state">Loading traces…</div> : null}

        {traces.isError ? (
          <div className="state state-error">
            <strong>Could not load traces</strong>
            <span>Check whether the API and database are running.</span>
          </div>
        ) : null}

        {traces.data && traces.data.traces.length === 0 ? (
          <div className="state">
            <strong>No traces found</strong>
            <span>Create a trace through the API to see it here.</span>
          </div>
        ) : null}

        {traces.data && traces.data.traces.length > 0 ? (
          <TraceTable traces={traces.data.traces} />
        ) : null}
      </section>
    </main>
  );
}
