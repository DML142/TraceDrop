export function MetadataBlock({
  title,
  metadata,
  compact = false,
}: {
  title: string;
  metadata: Record<string, unknown>;
  compact?: boolean;
}) {
  const empty = Object.keys(metadata).length === 0;

  return (
    <details className={compact ? "metadata metadata-compact" : "metadata"} open={!compact}>
      <summary>{title}</summary>
      {empty ? (
        <div className="metadata-empty">No metadata</div>
      ) : (
        <pre>{JSON.stringify(metadata, null, 2)}</pre>
      )}
    </details>
  );
}
