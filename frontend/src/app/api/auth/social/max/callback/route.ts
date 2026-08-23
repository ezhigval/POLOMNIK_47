import { NextRequest, NextResponse } from "next/server";
import { exchangeMax, redirectAfterOAuth } from "@/lib/auth/social-oauth";

export async function GET(request: NextRequest) {
  const code = request.nextUrl.searchParams.get("code");
  if (!code) {
    return NextResponse.redirect(new URL("/account/login?error=oauth_cancelled", request.url));
  }
  const result = await exchangeMax(code, request);
  if ("error" in result) {
    return NextResponse.redirect(
      new URL(`/account/login?error=${encodeURIComponent(result.error)}`, request.url),
    );
  }
  return redirectAfterOAuth(request, result);
}
