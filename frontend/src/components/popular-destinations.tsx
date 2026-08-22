import Link from "next/link";
import { TourCard } from "@/components/tour-card";
import { getPopularTours, type Tour } from "@/lib/api/tours";

export async function PopularDestinations() {
  let tours: Tour[] = [];
  try {
    tours = (await getPopularTours(8)) ?? [];
  } catch {
    return null;
  }

  if (tours.length === 0) {
    return null;
  }

  return (
    <section id="tours" className="scroll-mt-24 space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h2 className="font-display text-2xl font-semibold text-stone-900 sm:text-3xl">
            Популярные направления
          </h2>
        </div>
        <Link href="/search" className="btn-secondary text-sm">
          Все туры
        </Link>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {tours.map((tour) => (
          <TourCard key={tour.id} tour={tour} featured />
        ))}
      </div>
    </section>
  );
}

export function PopularDestinationsSkeleton() {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3" aria-busy="true">
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className="h-64 animate-pulse rounded-2xl bg-stone-200" />
      ))}
    </div>
  );
}
