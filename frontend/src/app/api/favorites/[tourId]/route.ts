import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { AUTH_COOKIE } from "@/lib/auth/session";

type RouteContext = { params: Promise<{ tourId: string }> };

async function proxy(method: string, tourId: string) {
  const cookieStore = await cookies();
  const token = cookieStore.get(AUTH_COOKIE)?.value;
  if (!token) {
    return NextResponse.json({ error: "Нужно войти в аккаунт" }, { status: 401 });
  }

  const response = await fetch(`${getApiBaseUrl()}/me/favorites/${tourId}`, {
    method,
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  const text = await response.text();
  return new NextResponse(text, { status: response.status, headers: { "Content-Type": "application/json" } });
}

export async function POST(_request: Request, context: RouteContext) {
  const { tourId } = await context.params;
  return proxy("POST", tourId);
}

export async function DELETE(_request: Request, context: RouteContext) {
  const { tourId } = await context.params;
  return proxy("DELETE", tourId);
}
