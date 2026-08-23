import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { ApiError, apiUrl, requestJson } from "@/lib/api/client";
import { AUTH_COOKIE } from "@/lib/auth/session";
import type { Passenger } from "@/lib/api/auth";

async function tokenFromCookie(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(AUTH_COOKIE)?.value ?? null;
}

export async function POST(request: Request) {
  const token = await tokenFromCookie();
  if (!token) {
    return NextResponse.json({ error: "Нужно войти в аккаунт" }, { status: 401 });
  }

  try {
    const body = await request.json();
    const passenger = await requestJson<{ data: Passenger }>(apiUrl("/me/passengers"), {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
      cache: "no-store",
    });
    return NextResponse.json(passenger, { status: 201 });
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
