import Link from "next/link";
import { FavoriteButton } from "@/components/favorite-button";
import { SlotsBadge } from "@/components/slots-badge";
import { TourImage } from "@/components/tour-image";
import {
  formatDateRange,
  formatPrice,
  formatTourDuration,
  getSlotsAvailability,
} from "@/lib/format";
import type { Tour } from "@/lib/api/tours";

type TourCardProps = {
  tour: Tour;
  featured?: boolean;
};

export function TourCard({ tour, featured = false }: TourCardProps) {
  const soldOut = getSlotsAvailability(tour.slots_left) === "sold_out";
  const duration = formatTourDuration(tour.date_start, tour.date_end);

  return (
    <article
      className={`group flex flex-col overflow-hidden rounded-2xl border bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-lg ${
        featured
          ? "border-amber-200/80 ring-1 ring-amber-100"
          : "border-stone-200/80 hover:border-brand-200"
      }`}
    >
      <div className="relative">
        <Link href={`/tours/${tour.id}`} className="relative block">
          <TourImage tour={tour} overlay className="aspect-[16/10] w-full" />
          <div className="absolute left-3 top-3 flex flex-wrap gap-2">
          {tour.is_hot ? (
            <span className="rounded-full bg-amber-400/95 px-2.5 py-0.5 text-xs font-semibold text-amber-950 shadow-sm">
              Хит сезона
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
        <div className="mb-3 flex items-start justify-between gap-3">
          <div className="min-w-0">
            <Link href={`/tours/${tour.id}`}>
              <h2 className="font-display text-xl font-semibold leading-snug text-stone-900 transition group-hover:text-brand-800">
                {tour.title}
              </h2>
            </Link>
            {tour.location ? (
              <p className="mt-1 text-sm text-stone-500">{tour.location}</p>
            ) : null}
          </div>
        </div>

        <p className="mb-4 line-clamp-2 flex-1 text-sm leading-6 text-stone-600">
          {tour.description?.split("\n")[0] || "Описание появится позже."}
        </p>

        <div className="mb-4 space-y-2 border-t border-stone-100 pt-4">
          <p className="text-sm text-stone-600">{formatDateRange(tour.date_start, tour.date_end)}</p>
          <div className="flex flex-wrap items-center justify-between gap-2">
            <p className="text-xl font-semibold text-stone-900">
              {formatPrice(tour.price, tour.currency)}
              <span className="text-sm font-normal text-stone-500"> / чел.</span>
            </p>
            <SlotsBadge slotsLeft={tour.slots_left} />
          </div>
        </div>

        <Link
          href={`/tours/${tour.id}`}
          className={`btn-primary w-full text-center ${soldOut ? "pointer-events-none opacity-50" : ""}`}
          aria-disabled={soldOut}
        >
          {soldOut ? "Мест нет" : "Подробнее и заявка"}
        </Link>
      </div>
    </article>
  );
}
