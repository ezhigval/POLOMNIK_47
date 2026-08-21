import { NextResponse } from "next/server";
import { safeReturnUrl } from "@/lib/site-nav";
import { siteConfig } from "@/lib/site-config";

function googleEnabled() {
  return Boolean(process.env.GOOGLE_CLIENT_ID && process.env.GOOGLE_CLIENT_SECRET);
}

function encodeOAuthState(returnUrl: string): string {
  return Buffer.from(returnUrl, "utf8").toString("base64url");
}

export async function GET(request: Request) {
  if (!googleEnabled()) {
    return NextResponse.json({ error: "Google OAuth is not configured" }, { status: 503 });
  }

  const requestUrl = new URL(request.url);
  const returnUrl = safeReturnUrl(requestUrl.searchParams.get("returnUrl"));
  const redirectUri = `${siteConfig.url.replace(/\/$/, "")}/api/auth/google/callback`;
  const params = new URLSearchParams({
    client_id: process.env.GOOGLE_CLIENT_ID!,
    redirect_uri: redirectUri,
    response_type: "code",
    scope: "openid email profile",
    access_type: "online",
    prompt: "select_account",
    state: encodeOAuthState(returnUrl),
  });

  return NextResponse.redirect(`https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`);
}
