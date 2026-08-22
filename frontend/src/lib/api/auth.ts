import { ApiError, apiUrl, requestJson } from "@/lib/api/client";

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

async function authRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const body = await requestJson<{ data?: T }>(apiUrl(path), {
    cache: "no-store",
    ...init,
  });
  if (body?.data === undefined) {
    throw new ApiError(500, "INVALID_RESPONSE", "Некорректный ответ сервера");
  }
  return body.data;
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
  return authRequest<MyBooking[]>("/me/bookings", {
    headers: { Authorization: `Bearer ${token}` },
  });
}
