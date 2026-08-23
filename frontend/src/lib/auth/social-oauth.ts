import { createHash, createHmac, timingSafeEqual } from "crypto";
import { NextRequest, NextResponse } from "next/server";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { AUTH_COOKIE } from "@/lib/auth/session";
import { SOCIAL_AUTH_PATHS, socialRedirectURI } from "@/lib/auth/social-paths";

export const UNAVAILABLE = "Пока что недоступно, используйте другой вариант.";

function siteOrigin(request: NextRequest): string {
  return process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, "") || new URL(request.url).origin;
}

async function completeOAuth(
  input: {
    provider: string;
    subject: string;
    email?: string;
    name?: string;
    phone?: string;
  },
  request: NextRequest,
): Promise<OAuthCompleteResult | { error: string; status: number }> {
  const secret = process.env.INTERNAL_API_SECRET;
  if (!secret) {
    return { error: UNAVAILABLE, status: 503 };
  }
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    "X-Internal-Secret": secret,
  };
  const session = request.cookies.get(AUTH_COOKIE)?.value;
  if (session) {
    headers.Authorization = `Bearer ${session}`;
  }
  const response = await fetch(`${getApiBaseUrl()}/auth/oauth`, {
    method: "POST",
    headers,
    body: JSON.stringify({
      provider: input.provider,
      subject: input.subject,
      email: input.email ?? "",
      name: input.name ?? "",
      phone: input.phone ?? "",
    }),
    cache: "no-store",
  }).catch(() => null);
  if (!response) {
    return { error: "Сервер авторизации недоступен", status: 503 };
  }
  const body = await response.json().catch(() => null);
  if (!response.ok) {
    return {
      error: body?.error?.message ?? UNAVAILABLE,
      status: response.status,
    };
  }
  const token = body?.data?.token;
  if (typeof token !== "string" || !token) {
    return { error: "Пустой ответ OAuth", status: 502 };
  }
  const kept = Array.isArray(body?.data?.kept_fields)
    ? (body.data.kept_fields as unknown[]).filter((item): item is string => typeof item === "string")
    : [];
  return {
    token,
    linked: Boolean(body?.data?.linked),
    merged: Boolean(body?.data?.merged),
    kept_fields: kept,
  };
}

export type OAuthCompleteResult = {
  token: string;
  linked?: boolean;
  merged?: boolean;
  kept_fields?: string[];
};

const KEPT_FIELDS = new Set(["name", "email", "phone"]);

export function redirectAfterOAuth(request: NextRequest, result: OAuthCompleteResult): NextResponse {
  const target = result.linked
    ? accountLinkURL(request, result)
    : new URL("/account/trips", request.url);
  const response = NextResponse.redirect(target);
  setSessionCookie(response, result.token);
  return response;
}

function accountLinkURL(request: NextRequest, result: OAuthCompleteResult): URL {
  const url = new URL("/account", request.url);
  url.searchParams.set("linked", "1");
  if (result.merged) {
    url.searchParams.set("merged", "1");
    const kept = (result.kept_fields ?? []).filter((field) => KEPT_FIELDS.has(field)).join(",");
    if (kept) {
      url.searchParams.set("kept", kept);
    }
  }
  return url;
}

export function setSessionCookie(response: NextResponse, token: string) {
  response.cookies.set(AUTH_COOKIE, token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24 * 7,
  });
}

export function oauthStartURL(
  provider: "yandex" | "vk" | "max",
  request: NextRequest,
): string | null {
  const origin = siteOrigin(request);
  if (provider === "yandex") {
    const id = process.env.YANDEX_OAUTH_CLIENT_ID;
    if (!id) return null;
    const redirect = socialRedirectURI("yandexCallback", origin);
    return `https://oauth.yandex.ru/authorize?response_type=code&client_id=${encodeURIComponent(id)}&redirect_uri=${encodeURIComponent(redirect)}`;
  }
  if (provider === "vk") {
    const id = process.env.VK_OAUTH_CLIENT_ID;
    if (!id) return null;
    const redirect = socialRedirectURI("vkCallback", origin);
    return `https://oauth.vk.com/authorize?response_type=code&client_id=${encodeURIComponent(id)}&redirect_uri=${encodeURIComponent(redirect)}&scope=email`;
  }
  const id = process.env.MAX_OAUTH_CLIENT_ID;
  const authorize = process.env.MAX_OAUTH_AUTHORIZE_URL;
  if (!id || !authorize) return null;
  const redirect = socialRedirectURI("maxCallback", origin);
  return `${authorize}?response_type=code&client_id=${encodeURIComponent(id)}&redirect_uri=${encodeURIComponent(redirect)}`;
}

export async function exchangeYandex(code: string, request: NextRequest) {
  const id = process.env.YANDEX_OAUTH_CLIENT_ID;
  const secret = process.env.YANDEX_OAUTH_CLIENT_SECRET;
  if (!id || !secret) return { error: UNAVAILABLE, status: 503 as const };
  const redirect = socialRedirectURI("yandexCallback", siteOrigin(request));
  const tokenRes = await fetch("https://oauth.yandex.ru/token", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      grant_type: "authorization_code",
      code,
      client_id: id,
      client_secret: secret,
      redirect_uri: redirect,
    }),
  });
  const tokenBody = await tokenRes.json().catch(() => null);
  if (!tokenRes.ok || !tokenBody?.access_token) {
    return { error: "Не удалось получить токен Яндекс", status: 502 as const };
  }
  const infoRes = await fetch("https://login.yandex.ru/info?format=json", {
    headers: { Authorization: `OAuth ${tokenBody.access_token}` },
  });
  const info = await infoRes.json().catch(() => null);
  if (!infoRes.ok || !info?.id) {
    return { error: "Не удалось получить профиль Яндекс", status: 502 as const };
  }
  return completeOAuth({
    provider: "yandex",
    subject: String(info.id),
    email: info.default_email || info.emails?.[0] || "",
    name: info.real_name || info.display_name || info.login || "",
  }, request);
}

export async function exchangeVK(code: string, request: NextRequest) {
  const id = process.env.VK_OAUTH_CLIENT_ID;
  const secret = process.env.VK_OAUTH_CLIENT_SECRET;
  if (!id || !secret) return { error: UNAVAILABLE, status: 503 as const };
  const redirect = socialRedirectURI("vkCallback", siteOrigin(request));
  const tokenRes = await fetch("https://oauth.vk.com/access_token", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      client_id: id,
      client_secret: secret,
      redirect_uri: redirect,
      code,
    }),
  });
  const tokenBody = await tokenRes.json().catch(() => null);
  if (!tokenRes.ok || !tokenBody?.user_id) {
    return { error: "Не удалось получить токен VK", status: 502 as const };
  }
  return completeOAuth({
    provider: "vk",
    subject: String(tokenBody.user_id),
    email: tokenBody.email || "",
    name: "",
  }, request);
}

export async function exchangeMax(code: string, request: NextRequest) {
  const id = process.env.MAX_OAUTH_CLIENT_ID;
  const secret = process.env.MAX_OAUTH_CLIENT_SECRET;
  const tokenURL = process.env.MAX_OAUTH_TOKEN_URL;
  const userInfoURL = process.env.MAX_OAUTH_USERINFO_URL;
  if (!id || !secret || !tokenURL || !userInfoURL) {
    return { error: UNAVAILABLE, status: 503 as const };
  }
  const redirect = socialRedirectURI("maxCallback", siteOrigin(request));
  const tokenRes = await fetch(tokenURL, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      grant_type: "authorization_code",
      code,
      client_id: id,
      client_secret: secret,
      redirect_uri: redirect,
    }),
  });
  const tokenBody = await tokenRes.json().catch(() => null);
  if (!tokenRes.ok || !tokenBody?.access_token) {
    return { error: "Не удалось получить токен Max", status: 502 as const };
  }
  const infoRes = await fetch(userInfoURL, {
    headers: { Authorization: `Bearer ${tokenBody.access_token}` },
  });
  const info = await infoRes.json().catch(() => null);
  const subject = info?.id || info?.sub;
  if (!infoRes.ok || !subject) {
    return { error: "Не удалось получить профиль Max", status: 502 as const };
  }
  return completeOAuth({
    provider: "max",
    subject: String(subject),
    email: info.email || "",
    name: info.name || info.display_name || "",
  }, request);
}

export function verifyTelegramLogin(payload: Record<string, string>): boolean {
  const botToken =
    process.env.TELEGRAM_LOGIN_BOT_TOKEN || process.env.TELEGRAM_BOT_TOKEN || "";
  if (!botToken || !payload.hash) {
    return false;
  }
  const { hash, ...rest } = payload;
  const checkString = Object.keys(rest)
    .sort()
    .map((key) => `${key}=${rest[key]}`)
    .join("\n");
  const secret = createHash("sha256").update(botToken).digest();
  const computed = createHmac("sha256", secret).update(checkString).digest("hex");
  try {
    return timingSafeEqual(Buffer.from(computed, "hex"), Buffer.from(hash, "hex"));
  } catch {
    return false;
  }
}

export async function completeTelegramOAuth(
  params: Record<string, string>,
  request: NextRequest,
): Promise<OAuthCompleteResult | { error: string; status: number }> {
  return completeOAuth(
    {
      provider: "telegram",
      subject: params.id,
      name: [params.first_name, params.last_name].filter(Boolean).join(" "),
      email: "",
      phone: "",
    },
    request,
  );
}

export { SOCIAL_AUTH_PATHS };
