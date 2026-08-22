import { ApiError, apiGet, apiGetList, requestJson, type DataEnvelope } from "./client";

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
  company_reply?: string;
  company_replied_at?: string | null;
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
  const body = await requestJson<DataEnvelope<CreateBookingResult>>("/api/bookings", {
    method: "POST",
    body: JSON.stringify(input),
  });
  if (!body?.data) {
    throw new ApiError(500, "INVALID_RESPONSE", "Некорректный ответ сервера");
  }
  return body.data;
}
