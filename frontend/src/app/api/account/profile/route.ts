import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { ApiError } from "@/lib/api/client";
import { updateUserProfile } from "@/lib/api/auth";
import { AUTH_COOKIE } from "@/lib/auth/session";

export async function PATCH(request: Request) {
  const cookieStore = await cookies();
  const token = cookieStore.get(AUTH_COOKIE)?.value;
  if (!token) {
    return NextResponse.json({ error: "Нужно войти в аккаунт" }, { status: 401 });
  }

  try {
    const body = await request.json();
    const user = await updateUserProfile(token, {
      name: body.name,
      email: body.email,
      phone: body.phone,
      phone_check_id: body.phone_check_id,
      website: body.website,
    });
    return NextResponse.json({ user });
  } catch (error) {
    if (error instanceof ApiError) {
      return NextResponse.json({ error: error.message }, { status: error.status });
    }
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Не удалось сохранить профиль" },
      { status: 400 },
    );
  }
}
