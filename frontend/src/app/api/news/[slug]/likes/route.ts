import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import { randomUUID } from "crypto";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { VISITOR_COOKIE, visitorCookieOptions } from "@/lib/visitor-id";

async function ensureVisitorId() {
  const cookieStore = await cookies();
  let visitorId = cookieStore.get(VISITOR_COOKIE)?.value;
  const created = !visitorId;
  if (!visitorId) {
    visitorId = randomUUID();
  }
  return { visitorId, created };
}

function withVisitorCookie(response: NextResponse, visitorId: string, created: boolean) {
  if (created) {
    response.cookies.set(VISITOR_COOKIE, visitorId, visitorCookieOptions());
  }
  return response;
}

export async function GET(_request: Request, context: { params: Promise<{ slug: string }> }) {
  const { slug } = await context.params;
  const { visitorId, created } = await ensureVisitorId();

  const upstream = await fetch(`${getApiBaseUrl()}/news/${encodeURIComponent(slug)}/likes`, {
    headers: { Cookie: `${VISITOR_COOKIE}=${visitorId}` },
    cache: "no-store",
  });
  const text = await upstream.text();
  return withVisitorCookie(
    new NextResponse(text, { status: upstream.status, headers: { "Content-Type": "application/json" } }),
    visitorId,
    created,
  );
}

export async function POST(_request: Request, context: { params: Promise<{ slug: string }> }) {
  const { slug } = await context.params;
  const { visitorId, created } = await ensureVisitorId();

  const upstream = await fetch(`${getApiBaseUrl()}/news/${encodeURIComponent(slug)}/likes`, {
    method: "POST",
    headers: { Cookie: `${VISITOR_COOKIE}=${visitorId}` },
    cache: "no-store",
  });
  const text = await upstream.text();
  return withVisitorCookie(
    new NextResponse(text, { status: upstream.status, headers: { "Content-Type": "application/json" } }),
    visitorId,
    created,
  );
}
