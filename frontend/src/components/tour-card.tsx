import Link from "next/link";
import { FavoriteButton } from "@/components/favorite-button";
import { SlotsBadge } from "@/components/slots-badge";
import { TourImage } from "@/components/tour-image";
import {
  formatDateRange,
  formatPrice,
  formatTourDuration,
} from "@/lib/format";
import { isRegularTour, isTourSoldOut, tourShowsPrice, type Tour } from "@/lib/api/tours";
import { tourPath } from "@/lib/tour-path";

type TourCardProps = {
  tour: Tour;
  featured?: boolean;
};

export function TourCard({ tour, featured = false }: TourCardProps) {
  const soldOut = isTourSoldOut(tour);
  const regular = isRegularTour(tour);
  const showPrice = tourShowsPrice(tour);
  const duration = regular ? "" : formatTourDuration(tour.date_start, tour.date_end);
  const highlighted = featured || tour.is_hot;

  return (
    <article
      className={`group flex flex-col overflow-hidden rounded-2xl border bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-lg ${
        highlighted
          ? "border-amber-200/80 ring-1 ring-amber-100"
          : "border-stone-200/80 hover:border-brand-200"
      }`}
    >
      <div className="relative">
        <Link href={tourPath(tour)} className="relative block">
          <TourImage tour={tour} overlay className="aspect-[16/10] w-full" />
          <div className="absolute left-3 top-3 flex flex-wrap gap-2">
            {tour.is_hot ? (
              <span className="rounded-full bg-amber-400/95 px-2.5 py-0.5 text-xs font-semibold text-amber-950 shadow-sm">
                Популярный
              </span>
            ) : null}
            {regular ? (
              <span className="rounded-full bg-white/90 px-2.5 py-0.5 text-xs font-medium text-stone-800 shadow-sm backdrop-blur-sm">
                Регулярный тур
              </span>
            ) : null}
            {duration ? (
              <span className="rounded-full bg-white/90 px-2.5 py-0.5 text-xs font-medium text-stone-800 shadow-sm backdrop-blur-sm">
                {duration}
              </span>
            ) : null}
          </div>
        </Link>
        <div className="absolute right-3 top-3 z-10">
          <FavoriteButton tourId={tour.id} compact />
        </div>
      </div>

      <div className="flex flex-1 flex-col p-5">
        <Link href={tourPath(tour)}>
          <h2 className="font-display text-xl font-semibold leading-snug text-stone-900 transition group-hover:text-brand-800">
            {tour.title}
          </h2>
        </Link>
        {tour.location ? <p className="mt-1 text-sm text-stone-500">{tour.location}</p> : null}

        {regular && !showPrice ? (
          <div className="mt-4 flex flex-1 items-end justify-end border-t border-stone-100 pt-4">
            <SlotsBadge slotsLeft={tour.slots_left} />
          </div>
        ) : (
          <dl className="mt-4 grid flex-1 grid-cols-2 gap-3 border-t border-stone-100 pt-4 text-sm">
            {regular ? null : (
              <>
                <div>
                  <dt className="text-xs uppercase tracking-wide text-stone-400">Даты</dt>
                  <dd className="mt-1 text-stone-800">{formatDateRange(tour.date_start, tour.date_end)}</dd>
                </div>
                <div>
                  <dt className="text-xs uppercase tracking-wide text-stone-400">Длительность</dt>
                  <dd className="mt-1 text-stone-800">{duration || "—"}</dd>
                </div>
              </>
            )}
            {showPrice ? (
              <div className="col-span-2 flex flex-wrap items-end justify-between gap-2">
                <div>
                  <dt className="text-xs uppercase tracking-wide text-stone-400">Стоимость</dt>
                  <dd className="mt-1 text-xl font-semibold text-stone-900">
                    {formatPrice(tour.price, tour.currency)}
                    <span className="text-sm font-normal text-stone-500"> / чел.</span>
                  </dd>
                </div>
                <SlotsBadge slotsLeft={tour.slots_left} />
              </div>
            ) : (
              <div className="col-span-2 flex justify-end">
                <SlotsBadge slotsLeft={tour.slots_left} />
              </div>
            )}
          </dl>
        )}

        <Link
          href={tourPath(tour)}
          className={`btn-primary mt-4 w-full text-center ${soldOut ? "pointer-events-none opacity-50" : ""}`}
          aria-disabled={soldOut}
        >
          {soldOut ? "Мест нет" : "Смотреть и записаться"}
        </Link>
      </div>
    </article>
  );
}
