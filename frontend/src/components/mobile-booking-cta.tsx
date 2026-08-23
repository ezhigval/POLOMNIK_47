"use client";

import { formatDateRange, formatPrice, formatTourDuration } from "@/lib/format";

type MobileBookingCTAProps = {
  disabled?: boolean;
  price?: number;
  currency?: string;
  dateStart?: string;
  dateEnd?: string;
};

export function MobileBookingCTA({
  disabled = false,
  price,
  currency = "RUB",
  dateStart,
  dateEnd,
}: MobileBookingCTAProps) {
  if (disabled) return null;

  const duration = dateStart && dateEnd ? formatTourDuration(dateStart, dateEnd) : "";
  const dates = dateStart && dateEnd ? formatDateRange(dateStart, dateEnd) : "";

  return (
    <div className="fixed inset-x-0 bottom-0 z-30 border-t border-stone-200 bg-white/95 p-4 backdrop-blur-md lg:hidden">
      <div className="flex items-center gap-3">
        <div className="min-w-0 flex-1">
          {dates ? <p className="truncate text-xs text-stone-500">{dates}</p> : null}
          {price ? (
            <p className="truncate font-semibold text-stone-900">
              {formatPrice(price, currency)}
              <span className="text-sm font-normal text-stone-500"> / чел.</span>
              {duration ? <span className="ml-2 text-xs font-normal text-stone-500">{duration}</span> : null}
            </p>
          ) : null}
        </div>
        <a href="#booking-form" className="btn-primary shrink-0 px-6">
          Заявка
        </a>
      </div>
    </div>
  );
}
