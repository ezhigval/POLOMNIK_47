import { ApiError, apiGet, apiGetList } from "./client";

export type Tour = {
  id: string;
  slug: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  date_start: string;
  date_end: string;
  slots_total: number;
  slots_left: number;
  location: string;
  images: string[];
  is_hot: boolean;
};

export type Review = {
  id: string;
  tour_id: string;
  client_name: string;
  rating: number;
  text: string;
  created_at: string;
};

export type CreateBookingInput = {
  tour_id: string;
  name: string;
  phone: string;
  email?: string;
  people_count: number;
  comment?: string;
};

export type CreateBookingResult = {
  status: string;
  booking_id: string;
  booking_status: string;
  total_price: number;
  integration_status: string;
};

export function getTours(params?: Record<string, string>) {
  const query = params ? `?${new URLSearchParams(params)}` : "";
  return apiGetList<Tour>(`/tours${query}`);
}

export function getTour(id: string) {
  return apiGet<Tour>(`/tours/${id}`);
}

export function getPopularTours(limit = 10) {
  return apiGet<Tour[]>(`/tours/popular?limit=${limit}`);
}

export function getTourReviews(tourId: string, page = 1, limit = 20) {
  return apiGetList<Review>(
    `/tours/${tourId}/reviews?page=${page}&limit=${limit}`,
  );
}

export function getReviews(page = 1, limit = 20) {
  return apiGetList<Review>(`/reviews?page=${page}&limit=${limit}`);
}

export async function createBooking(input: CreateBookingInput) {
  const response = await fetch("/api/bookings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });

  const text = await response.text();
  let body: { data?: CreateBookingResult; error?: { message?: string } } | null = null;
  try {
    body = text ? JSON.parse(text) : null;
  } catch {
    body = null;
  }

  if (!response.ok) {
    const code = body?.error && "code" in body.error ? String((body.error as { code?: string }).code) : "UNKNOWN_ERROR";
    throw new ApiError(response.status, code, body?.error?.message ?? "Booking failed");
  }

  if (!body?.data) {
    throw new ApiError(500, "INVALID_RESPONSE", "Invalid booking response");
  }

  return body.data;
}
