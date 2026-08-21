"use client";

import { formatPrice } from "@/lib/format";

type MobileBookingCTAProps = {
  disabled?: boolean;
  price?: number;
  currency?: string;
};

export function MobileBookingCTA({
  disabled = false,
  price,
  currency = "RUB",
}: MobileBookingCTAProps) {
  if (disabled) return null;

  return (
    <div className="fixed inset-x-0 bottom-0 z-30 border-t border-stone-200 bg-white/95 p-4 backdrop-blur-md lg:hidden">
      <div className="flex items-center gap-3">
        {price ? (
          <div className="min-w-0 flex-1">
            <p className="text-xs text-stone-500">от</p>
            <p className="truncate font-semibold text-stone-900">
              {formatPrice(price, currency)}
              <span className="text-sm font-normal text-stone-500"> / чел.</span>
            </p>
          </div>
        ) : null}
        <a href="#booking-form" className="btn-primary shrink-0 px-6">
          Заявка
        </a>
      </div>
    </div>
  );
}
