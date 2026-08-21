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

export async function isAdminAuthenticated(): Promise<boolean> {
  const adminToken = process.env.ADMIN_TOKEN;
  if (!adminToken) {
    return false;
  }

  const cookieStore = await cookies();
  const session = cookieStore.get(ADMIN_SESSION_COOKIE)?.value;
  if (!session) {
    return false;
  }

  return verifyAdminSessionValue(session, adminToken);
}
