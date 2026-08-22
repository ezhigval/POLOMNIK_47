import { NextResponse } from "next/server";
import { registerUser } from "@/lib/api/auth";
import { AUTH_COOKIE } from "@/lib/auth/session";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const result = await registerUser({
      name: body.name,
      email: body.email,
      phone: body.phone,
      password: body.password,
    });

    const response = NextResponse.json({ user: result.user });
    response.cookies.set(AUTH_COOKIE, result.token, {
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      path: "/",
      maxAge: 60 * 60 * 24 * 7,
    });
    return response;
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Не удалось зарегистрироваться" },
      { status: 400 },
    );
  }
}
