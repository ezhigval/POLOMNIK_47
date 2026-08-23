import { NextResponse } from "next/server";
import { ApiError, apiUrl, requestJson } from "@/lib/api/client";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const result = await requestJson<{ data?: { message?: string } }>(apiUrl("/auth/reset-password"), {
      method: "POST",
      cache: "no-store",
      body: JSON.stringify({
        token: body.token,
        new_password: body.new_password,
      }),
    });
    return NextResponse.json({ message: result.data?.message ?? "OK" });
  } catch (error) {
    if (error instanceof ApiError) {
      return NextResponse.json({ error: error.message }, { status: error.status });
    }
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Не удалось обновить пароль" },
      { status: 500 },
    );
  }
}
