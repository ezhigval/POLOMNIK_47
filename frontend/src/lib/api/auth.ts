import { getApiBaseUrl } from "@/lib/api/base-url";

export type User = {
  id: string;
  email: string;
  phone: string;
  name: string;
  created_at: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

export type MyBooking = {
  id: string;
  tour_id: string;
  name: string;
  phone: string;
  email: string;
  people_count: number;
  status: string;
  total_price: number;
  comment: string;
  created_at: string;
};

function apiUrl(path: string): string {
  return `${getApiBaseUrl()}${path}`;
}

async function authRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(apiUrl(path), {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
    cache: "no-store",
  });

  const text = await response.text();
  let body: { data?: T; error?: { message?: string } } | null = null;
  try {
    body = text ? JSON.parse(text) : null;
  } catch {
    body = null;
  }

  if (!response.ok) {
    throw new Error(body?.error?.message ?? "Request failed");
  }

  return body?.data as T;
}

export async function registerUser(input: {
  name: string;
  email: string;
  phone: string;
  password: string;
}): Promise<AuthResponse> {
  return authRequest<AuthResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function loginUser(input: {
  login: string;
  password: string;
}): Promise<AuthResponse> {
  return authRequest<AuthResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function fetchCurrentUser(token: string): Promise<User> {
  return authRequest<User>("/me", {
    headers: { Authorization: `Bearer ${token}` },
  });
}

export async function fetchMyBookings(token: string): Promise<MyBooking[]> {
  const response = await fetch(apiUrl("/me/bookings"), {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  const body = await response.json();
  if (!response.ok) {
    throw new Error(body?.error?.message ?? "Failed to load bookings");
  }
  return body.data as MyBooking[];
}
