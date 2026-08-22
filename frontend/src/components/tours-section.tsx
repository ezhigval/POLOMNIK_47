import { Suspense } from "react";
import { Pagination } from "@/components/pagination";
import { TourCard } from "@/components/tour-card";
import { ActiveFilterChips, TourFilters } from "@/components/tour-filters";
import { ApiError } from "@/lib/api/client";
import { getPopularTours, getTours } from "@/lib/api/tours";
import {
  hasActiveFilters,
  parseTourFilters,
  toTourQueryParams,
  type TourFilterValues,
} from "@/lib/tour-filters";

type ToursSectionProps = {
  searchParams: Record<string, string | string[] | undefined>;
  basePath?: string;
  showPopular?: boolean;
  showFilters?: boolean;
};

async function loadToursData(filters: TourFilterValues, showPopularBlock: boolean) {
  const query = toTourQueryParams(filters);
  const showPopular = showPopularBlock && !hasActiveFilters(filters);

  try {
    const [toursResponse, popular] = await Promise.all([
      getTours(query),
      showPopular ? getPopularTours(6) : Promise.resolve([]),
    ]);

    return {
      error: null as string | null,
      toursResponse,
      popular,
      showPopular,
    };
  } catch (err) {
    return {
      error:
        err instanceof ApiError
          ? err.message
          : "Не удалось загрузить туры. Попробуйте обновить страницу.",
      toursResponse: null,
      popular: [],
      showPopular,
    };
  }
}

async function ToursContent({
  filters,
  showPopularBlock,
  basePath,
}: {
  filters: TourFilterValues;
  showPopularBlock: boolean;
  basePath: string;
}) {
  const { error, toursResponse, popular, showPopular } = await loadToursData(filters, showPopularBlock);

  if (error) {
    return (
      <div className="rounded-2xl border border-red-200 bg-red-50 p-5 text-red-800" role="alert">
        {error}
      </div>
    );
  }

  if (!toursResponse) {
    return null;
  }

  return (
    <>
      {showPopular && popular.length > 0 ? (
        <section className="mb-10">
          <h2 className="mb-4 text-xl font-semibold">Популярные туры</h2>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {popular.map((tour) => (
              <TourCard key={tour.id} tour={tour} featured />
            ))}
          </div>
        </section>
      ) : null}

      <section>
        <div className="mb-4 flex items-center justify-between gap-3">
          <h2 className="text-xl font-semibold">
            {hasActiveFilters(filters) ? "Результаты поиска" : "Все туры"}
          </h2>
          <span className="text-sm text-stone-500">Найдено: {toursResponse.meta.total}</span>
        </div>

        {toursResponse.data.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-stone-300 bg-white p-8 text-center text-stone-500">
            {hasActiveFilters(filters)
              ? "По выбранным фильтрам туры не найдены."
              : "Пока нет доступных туров."}
          </div>
        ) : (
          <>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {toursResponse.data.map((tour) => (
                <TourCard key={tour.id} tour={tour} />
              ))}
            </div>
            <Pagination meta={toursResponse.meta} filters={filters} basePath={basePath} />
          </>
        )}
      </section>
    </>
  );
}

export function ToursSection({
  searchParams,
  basePath = "/search",
  showPopular = true,
  showFilters = true,
}: ToursSectionProps) {
  const filters = parseTourFilters(searchParams);

  return (
    <>
      {showFilters ? (
        <Suspense fallback={<div className="h-48 animate-pulse rounded-2xl bg-stone-200" />}>
          <TourFilters initial={filters} basePath={basePath} />
        </Suspense>
      ) : null}
      <div className={showFilters ? "mt-4" : undefined}>
        <ActiveFilterChips filters={filters} basePath={basePath} />
      </div>
      <div className="mt-6">
        <Suspense fallback={<ToursSkeleton />}>
          <ToursContent filters={filters} showPopularBlock={showPopular} basePath={basePath} />
        </Suspense>
      </div>
    </>
  );
}

function ToursSkeleton() {
  return (
    <div className="space-y-4" aria-live="polite" aria-busy="true">
      <div className="h-6 w-40 animate-pulse rounded bg-stone-200" />
      <div className="grid gap-4 md:grid-cols-2">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="h-48 animate-pulse rounded-2xl bg-stone-200" />
        ))}
      </div>
    </div>
  );
}
