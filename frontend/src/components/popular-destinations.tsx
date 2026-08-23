import Link from "next/link";
import { TourSchedule } from "@/components/tour-schedule";
import { getTours, type Tour } from "@/lib/api/tours";

export async function PopularDestinations() {
  let tours: Tour[] = [];
  try {
    const response = await getTours({ limit: "8" });
    tours = response.data ?? [];
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
            Ближайшие выезды
          </h2>
          <p className="mt-2 text-sm text-stone-600">
            Даты, длительность и цена — из карточки тура, по дате выезда.
          </p>
        </div>
        <Link href="/search" className="btn-secondary text-sm">
          Всё расписание
        </Link>
      </div>

      <TourSchedule tours={tours} />
    </section>
  );
}

export function PopularDestinationsSkeleton() {
  return (
    <div className="h-64 animate-pulse rounded-2xl bg-stone-200" aria-busy="true" />
  );
}
