import { NextResponse } from "next/server";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { AUTH_COOKIE } from "@/lib/auth/session";
import { safeReturnUrl } from "@/lib/site-nav";
import { siteConfig } from "@/lib/site-config";

function decodeOAuthState(state: string | null): string {
  if (!state) {
    return "/account/trips";
  }
  try {
    return safeReturnUrl(Buffer.from(state, "base64url").toString("utf8"));
  } catch {
    return "/account/trips";
  }
}

function loginRedirect(error: string, returnUrl?: string) {
  const url = new URL("/account/login", siteConfig.url);
  url.searchParams.set("error", error);
  if (returnUrl && returnUrl !== "/account/trips") {
    url.searchParams.set("returnUrl", returnUrl);
  }
  return NextResponse.redirect(url);
}

export async function GET(request: Request) {
  const clientId = process.env.GOOGLE_CLIENT_ID;
  const clientSecret = process.env.GOOGLE_CLIENT_SECRET;
  if (!clientId || !clientSecret) {
    return loginRedirect("oauth_not_configured");
  }

  const url = new URL(request.url);
  const returnUrl = decodeOAuthState(url.searchParams.get("state"));
  const code = url.searchParams.get("code");
  if (!code) {
    return loginRedirect("oauth_cancelled", returnUrl);
  }

  const redirectUri = `${siteConfig.url.replace(/\/$/, "")}/api/auth/google/callback`;
  const tokenResponse = await fetch("https://oauth2.googleapis.com/token", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      code,
      client_id: clientId,
      client_secret: clientSecret,
      redirect_uri: redirectUri,
      grant_type: "authorization_code",
    }),
  });

  const tokenPayload = await tokenResponse.json();
  if (!tokenResponse.ok || !tokenPayload.access_token) {
    return loginRedirect("oauth_failed", returnUrl);
  }

  const profileResponse = await fetch("https://www.googleapis.com/oauth2/v3/userinfo", {
    headers: { Authorization: `Bearer ${tokenPayload.access_token}` },
  });
  const profile = await profileResponse.json();
  if (!profileResponse.ok || !profile.sub) {
    return loginRedirect("oauth_profile", returnUrl);
  }

  const backendResponse = await fetch(`${getApiBaseUrl()}/auth/oauth`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Internal-Secret": process.env.INTERNAL_API_SECRET ?? "",
    },
    body: JSON.stringify({
      provider: "google",
      subject: profile.sub,
      email: profile.email ?? "",
      name: profile.name ?? profile.email ?? "Пользователь",
    }),
    cache: "no-store",
  });

  const backendPayload = await backendResponse.json();
  if (!backendResponse.ok || !backendPayload?.data?.token) {
    return loginRedirect("oauth_backend", returnUrl);
  }

  const redirect = NextResponse.redirect(new URL(returnUrl, siteConfig.url));
  redirect.cookies.set(AUTH_COOKIE, backendPayload.data.token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24 * 7,
  });
  return redirect;
}
