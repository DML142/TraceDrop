"use client";

import Link from "next/link";

import type { Trace } from "@/types/trace";

function formatDuration(duration: number | null) {
  if (duration === null) {
    return "—";
  }

  if (duration < 1000) {
    return `${duration} ms`;
  }

  return `${(duration / 1000).toFixed(2)} s`;
}

function formatTime(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(new Date(value));
}

export function TraceTable({ traces }: { traces: Trace[] }) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Status</th>
            <th>Started</th>
            <th>Duration</th>
          </tr>
        </thead>
        <tbody>
          {traces.map((trace) => (
            <tr key={trace.id}>
              <td>
                <Link className="trace-link" href={`/traces/${trace.id}`}>
                  <div className="trace-name">{trace.name}</div>
                  <div className="trace-id">{trace.id}</div>
                </Link>
              </td>
              <td>
                <span className={`status status-${trace.status}`}>{trace.status}</span>
              </td>
              <td>{formatTime(trace.started_at)}</td>
              <td>{formatDuration(trace.duration_ms)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
