import type { User } from "@/lib/api/auth";
import { apiUrl, requestJson } from "@/lib/api/client";
import { AUTH_COOKIE } from "@/lib/auth/session";
import { cookies } from "next/headers";

async function authedFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
  const cookieStore = await cookies();
  const token = cookieStore.get(AUTH_COOKIE)?.value;
  if (!token) {
    throw new Error("Нужно войти в аккаунт");
  }

  const body = await requestJson<{ data?: T }>(apiUrl(path), {
    cache: "no-store",
    ...init,
    headers: {
      Authorization: `Bearer ${token}`,
      ...(init.headers ?? {}),
    },
  });
  return body.data as T;
}

export async function fetchFavoriteTourIds(): Promise<string[]> {
  return authedFetch("/me/favorites");
}

export type SupportMessage = {
  id: string;
  sender_type: "user" | "staff";
  body: string;
  created_at: string;
};

export type SupportThread = {
  id: string;
  subject: string;
  status: string;
  messages: SupportMessage[];
  updated_at: string;
};

export async function fetchSupportThread(): Promise<SupportThread> {
  return authedFetch("/me/support");
}

export async function sendSupportMessage(body: string): Promise<SupportThread> {
  return authedFetch("/me/support/messages", {
    method: "POST",
    body: JSON.stringify({ body }),
  });
}

export type BookingProfile = Pick<User, "name" | "phone" | "email">;

export function toBookingProfile(user: User | null): BookingProfile | null {
  if (!user) {
    return null;
  }
  return {
    name: user.name,
    phone: user.phone,
    email: user.email,
  };
}
