import { NextResponse } from "next/server";
import { ADMIN_SESSION_COOKIE } from "@/lib/auth/admin-session";
import { apiUrl } from "@/lib/api/client";

export async function POST(request: Request) {
  const body = await request.json().catch(() => null);
  const password = typeof body?.password === "string" ? body.password : "";
  const role = typeof body?.role === "string" ? body.role : "";
  if (!password) {
    return NextResponse.json({ error: "Введите пароль" }, { status: 400 });
  }

  const upstream = await fetch(apiUrl("/management/auth/login"), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ role, password }),
    cache: "no-store",
  }).catch(() => null);

  if (!upstream) {
    return NextResponse.json({ error: "API недоступен" }, { status: 503 });
  }

  const payload = await upstream.json().catch(() => null);
  if (!upstream.ok) {
    return NextResponse.json(
      { error: payload?.error?.message ?? "Неверный логин или пароль" },
      { status: upstream.status },
    );
  }

  const token = payload?.data?.token;
  if (typeof token !== "string" || !token) {
    return NextResponse.json({ error: "Пустой ответ авторизации" }, { status: 502 });
  }

  const response = NextResponse.json({
    ok: true,
    full_admin: Boolean(payload?.data?.full_admin),
    permissions: payload?.data?.permissions ?? [],
  });
  response.cookies.set(ADMIN_SESSION_COOKIE, token, {
    httpOnly: true,
    sameSite: "strict",
    secure: process.env.NODE_ENV === "production",
    path: "/management",
    maxAge: 60 * 60 * 8,
  });
  return response;
}
