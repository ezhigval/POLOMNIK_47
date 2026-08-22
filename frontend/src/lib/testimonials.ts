import { getReviews, getTours } from "@/lib/api/tours";
import { fallbackTestimonials } from "@/lib/site-content";
import type { Testimonial } from "@/components/testimonial-card";

export async function loadTestimonials(limit: number): Promise<Testimonial[]> {
  try {
    const [reviewsResponse, toursResponse] = await Promise.all([
      getReviews(1, Math.max(limit, 12)),
      getTours({ limit: "100" }),
    ]);

    if (reviewsResponse.data.length === 0) {
      return fallbackTestimonials.slice(0, limit);
    }

    const tourTitles = new Map(toursResponse.data.map((tour) => [tour.id, tour.title]));
    return reviewsResponse.data.slice(0, limit).map((review) => ({
      client_name: review.client_name,
      text: review.text,
      rating: review.rating,
      tour_title: tourTitles.get(review.tour_id) ?? "Паломнический тур",
      company_reply: review.company_reply,
    }));
  } catch {
    return fallbackTestimonials.slice(0, limit);
  }
}
