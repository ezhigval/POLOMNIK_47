import { NextRequest, NextResponse } from "next/server";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { setSessionCookie, UNAVAILABLE, verifyTelegramLogin } from "@/lib/auth/social-oauth";

export async function GET(request: NextRequest) {
  const params = Object.fromEntries(request.nextUrl.searchParams.entries());
  if (!params.id || !params.hash) {
    return NextResponse.redirect(new URL("/account/login?error=oauth_cancelled", request.url));
  }
  if (!verifyTelegramLogin(params)) {
    return NextResponse.redirect(new URL("/account/login?error=oauth_failed", request.url));
  }

  const secret = process.env.INTERNAL_API_SECRET;
  if (!secret) {
    return NextResponse.redirect(new URL(`/account/login?error=${encodeURIComponent(UNAVAILABLE)}`, request.url));
  }

  const response = await fetch(`${getApiBaseUrl()}/auth/oauth`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Internal-Secret": secret,
    },
    body: JSON.stringify({
      provider: "telegram",
      subject: params.id,
      name: [params.first_name, params.last_name].filter(Boolean).join(" "),
      email: "",
      phone: "",
    }),
    cache: "no-store",
  }).catch(() => null);

  const body = await response?.json().catch(() => null);
  const token = body?.data?.token;
  if (!response?.ok || typeof token !== "string") {
    return NextResponse.redirect(new URL("/account/login?error=oauth_backend", request.url));
  }

  const redirect = NextResponse.redirect(new URL("/account/trips", request.url));
  setSessionCookie(redirect, token);
  return redirect;
}
