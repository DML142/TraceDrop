import { TraceDetail } from "@/features/traces/trace-detail";

export default async function TracePage({
  params,
}: {
  params: Promise<{ traceId: string }>;
}) {
  const { traceId } = await params;

  return <TraceDetail traceId={traceId} />;
}
