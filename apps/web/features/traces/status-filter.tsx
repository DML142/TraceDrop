"use client";

import type { TraceStatus } from "@/types/trace";

const filters: Array<{ label: string; value?: TraceStatus }> = [
  { label: "All" },
  { label: "Running", value: "running" },
  { label: "Completed", value: "completed" },
  { label: "Failed", value: "failed" },
];

export function StatusFilter({
  value,
  onChange,
}: {
  value?: TraceStatus;
  onChange: (status?: TraceStatus) => void;
}) {
  return (
    <div className="filter-group" role="group" aria-label="Trace status">
      {filters.map((filter) => {
        const active = filter.value === value;

        return (
          <button
            className={active ? "filter-button filter-button-active" : "filter-button"}
            key={filter.label}
            onClick={() => onChange(filter.value)}
            type="button"
          >
            {filter.label}
          </button>
        );
      })}
    </div>
  );
}
