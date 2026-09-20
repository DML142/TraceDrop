import { NextRequest, NextResponse } from "next/server";

const apiUrl = process.env.API_URL ?? process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function GET(request: NextRequest) {
  const upstream = new URL("/api/v1/traces", apiUrl);
  const status = request.nextUrl.searchParams.get("status");

  if (status) {
    upstream.searchParams.set("status", status);
  }

  upstream.searchParams.set("limit", "50");

  try {
    const response = await fetch(upstream, {
      cache: "no-store",
    });

    const body = await response.text();

    return new NextResponse(body, {
      status: response.status,
      headers: {
        "content-type": response.headers.get("content-type") ?? "application/json",
      },
    });
  } catch {
    return NextResponse.json(
      {
        error: {
          code: "api_unavailable",
          message: "Trace API is unavailable",
        },
      },
      { status: 502 },
    );
  }
}
