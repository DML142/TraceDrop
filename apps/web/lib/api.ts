import type { TraceListResponse, TraceStatus } from "@/types/trace";

export async function getTraces(status?: TraceStatus): Promise<TraceListResponse> {
  const params = new URLSearchParams();
  if (status) {
    params.set("status", status);
  }

  const response = await fetch(`/api/traces?${params.toString()}`, {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Failed to load traces");
  }

  return response.json() as Promise<TraceListResponse>;
}
