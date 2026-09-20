import type { TraceEvent } from "@/types/trace";

import { MetadataBlock } from "./metadata-block";

function formatTimestamp(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    fractionalSecondDigits: 3,
  }).format(new Date(value));
}

function formatDuration(value: number | null) {
  if (value === null) {
    return null;
  }
  return value < 1000 ? `${value} ms` : `${(value / 1000).toFixed(2)} s`;
}

export function TraceTimeline({ events }: { events: TraceEvent[] }) {
  if (events.length === 0) {
    return (
      <div className="state">
        <strong>No events recorded</strong>
        <span>This trace exists, but its timeline is empty.</span>
      </div>
    );
  }

  return (
    <ol className="timeline">
      {events.map((event) => {
        const duration = formatDuration(event.duration_ms);

        return (
          <li className="timeline-item" key={event.id}>
            <div className={`timeline-marker timeline-marker-${event.event_type}`} />
            <div className="timeline-content">
              <div className="timeline-row">
                <div>
                  <span className="event-type">{event.event_type}</span>
                  {event.message ? <p className="event-message">{event.message}</p> : null}
                </div>
                <div className="timeline-time">
                  <span>{formatTimestamp(event.timestamp)}</span>
                  {duration ? <strong>{duration}</strong> : null}
                </div>
              </div>

              {Object.keys(event.metadata).length > 0 ? (
                <MetadataBlock title="Event metadata" metadata={event.metadata} compact />
              ) : null}
            </div>
          </li>
        );
      })}
    </ol>
  );
}
