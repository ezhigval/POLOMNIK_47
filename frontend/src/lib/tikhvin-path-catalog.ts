import type { Tour } from "@/lib/api/tours";
import { featuredRoute } from "@/lib/featured-route";
import type { TourFilterValues } from "@/lib/tour-filters";

const TIKHVIN_PATH_TITLE = /тихвинский\s+путь/i;

export function isTikhvinPathCatalog(filters: TourFilterValues): boolean {
  return filters.route === featuredRoute.id || filters.destination === "tikhvin";
}

export function isTikhvinPathTour(tour: Pick<Tour, "title" | "slug">): boolean {
  return tour.slug === featuredRoute.id || TIKHVIN_PATH_TITLE.test(tour.title);
}

export function partitionTikhvinPathTours(tours: Tour[]): { featured: Tour[]; others: Tour[] } {
  const featured: Tour[] = [];
  const others: Tour[] = [];
  for (const tour of tours) {
    if (isTikhvinPathTour(tour)) {
      featured.push(tour);
    } else {
      others.push(tour);
    }
  }
  return { featured, others };
}
