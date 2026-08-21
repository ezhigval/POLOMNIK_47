import { cookies } from "next/headers";
import { fetchCurrentUser, type User } from "@/lib/api/auth";

export const AUTH_COOKIE = "polomnik_token";

export async function getSessionUser(): Promise<User | null> {
  const cookieStore = await cookies();
  const token = cookieStore.get(AUTH_COOKIE)?.value;
  if (!token) {
    return null;
  }

  try {
    return await fetchCurrentUser(token);
  } catch {
    return null;
  }
}

export async function getAuthToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(AUTH_COOKIE)?.value ?? null;
}
