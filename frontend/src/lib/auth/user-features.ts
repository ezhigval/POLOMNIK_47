import type { User } from "@/lib/api/auth";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { AUTH_COOKIE } from "@/lib/auth/session";
import { cookies } from "next/headers";

async function authedFetch(path: string, init: RequestInit = {}) {
  const cookieStore = await cookies();
  const token = cookieStore.get(AUTH_COOKIE)?.value;
  if (!token) {
    throw new Error("Authentication required");
  }

  const response = await fetch(`${getApiBaseUrl()}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...(init.headers ?? {}),
    },
    cache: "no-store",
  });

  const text = await response.text();
  const body = text ? JSON.parse(text) : null;
  if (!response.ok) {
    throw new Error(body?.error?.message ?? "Request failed");
  }
  return body?.data;
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
