import { createHash, timingSafeEqual } from "crypto";
import { cookies } from "next/headers";

export const ADMIN_SESSION_COOKIE = "polomnik_admin_session";

export function adminSessionToken(adminToken: string): string {
  return createHash("sha256").update(`admin:${adminToken}`).digest("hex");
}

export function verifyAdminPassword(password: string, adminToken: string): boolean {
  if (!password || !adminToken) {
    return false;
  }

  const expected = Buffer.from(adminToken, "utf8");
  const actual = Buffer.from(password, "utf8");
  if (expected.length !== actual.length) {
    return false;
  }

  return timingSafeEqual(expected, actual);
}

export function verifyAdminSessionValue(sessionValue: string, adminToken: string): boolean {
  if (!sessionValue || !adminToken) {
    return false;
  }

  const expected = Buffer.from(adminSessionToken(adminToken), "utf8");
  const actual = Buffer.from(sessionValue, "utf8");
  if (expected.length !== actual.length) {
    return false;
  }

  return timingSafeEqual(expected, actual);
}

export function isManagementJwt(sessionValue: string): boolean {
  return sessionValue.split(".").length === 3;
}

export async function getAdminSessionCookie(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(ADMIN_SESSION_COOKIE)?.value ?? null;
}

export async function isAdminAuthenticated(): Promise<boolean> {
  const session = await getAdminSessionCookie();
  if (!session) {
    return false;
  }
  if (isManagementJwt(session)) {
    return true;
  }
  const adminToken = process.env.ADMIN_TOKEN;
  if (!adminToken) {
    return false;
  }
  return verifyAdminSessionValue(session, adminToken);
}
