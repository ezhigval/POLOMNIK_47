import { NextRequest, NextResponse } from "next/server";
import { oauthStartURL, UNAVAILABLE } from "@/lib/auth/social-oauth";

export async function GET(request: NextRequest) {
  const url = oauthStartURL("vk", request);
  if (!url) {
    return NextResponse.redirect(new URL(`/account/login?error=${encodeURIComponent(UNAVAILABLE)}`, request.url));
  }
  return NextResponse.redirect(url);
}
