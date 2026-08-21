import { cache } from "react";
import { getTour, getTourReviews } from "@/lib/api/tours";

export const getCachedTour = cache((id: string) => getTour(id));

export const getCachedTourReviews = cache((tourId: string, page = 1, limit = 20) =>
  getTourReviews(tourId, page, limit),
);
