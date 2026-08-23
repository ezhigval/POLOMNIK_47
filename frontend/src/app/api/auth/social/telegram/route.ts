import { NextRequest, NextResponse } from "next/server";
import { completeTelegramOAuth, redirectAfterOAuth, verifyTelegramLogin } from "@/lib/auth/social-oauth";

export async function GET(request: NextRequest) {
  const params = Object.fromEntries(request.nextUrl.searchParams.entries());
  if (!params.id || !params.hash) {
    return NextResponse.redirect(new URL("/account/login?error=oauth_cancelled", request.url));
  }
  if (!verifyTelegramLogin(params)) {
    return NextResponse.redirect(new URL("/account/login?error=oauth_failed", request.url));
  }

  const result = await completeTelegramOAuth(params, request);
  if ("error" in result) {
    return NextResponse.redirect(new URL("/account/login?error=oauth_backend", request.url));
  }
  return redirectAfterOAuth(request, result);
}
