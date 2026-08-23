import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { ApiError, apiUrl, requestJson } from "@/lib/api/client";
import { AUTH_COOKIE } from "@/lib/auth/session";
import type { Passenger } from "@/lib/api/auth";

type RouteContext = {
  params: Promise<{ id: string }>;
};

async function tokenFromCookie(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(AUTH_COOKIE)?.value ?? null;
}

export async function PATCH(request: Request, context: RouteContext) {
  const token = await tokenFromCookie();
  if (!token) {
    return NextResponse.json({ error: "Нужно войти в аккаунт" }, { status: 401 });
  }
  const { id } = await context.params;

  try {
    const body = await request.json();
    const passenger = await requestJson<{ data: Passenger }>(apiUrl(`/me/passengers/${id}`), {
      method: "PATCH",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
      cache: "no-store",
    });
    return NextResponse.json(passenger);
  } catch (error) {
    if (error instanceof ApiError) {
      return NextResponse.json({ error: error.message }, { status: error.status });
    }
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Не удалось сохранить пассажира" },
      { status: 400 },
    );
  }
}

export async function DELETE(_request: Request, context: RouteContext) {
  const token = await tokenFromCookie();
  if (!token) {
    return NextResponse.json({ error: "Нужно войти в аккаунт" }, { status: 401 });
  }
  const { id } = await context.params;

  try {
    await requestJson(apiUrl(`/me/passengers/${id}`), {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    });
    return NextResponse.json({ removed: true });
  } catch (error) {
    if (error instanceof ApiError) {
      return NextResponse.json({ error: error.message }, { status: error.status });
    }
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Не удалось удалить пассажира" },
      { status: 400 },
    );
  }
}
