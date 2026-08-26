import { formatPrice } from "@/lib/format";

type TourPriceProps = {
  price: number | null | undefined;
  originalPrice?: number | null;
  currency: string;
  className?: string;
  suffix?: string;
};

/** Renders unit price; strikethrough original when burning discount applies. */
export function TourPrice({
  price,
  originalPrice,
  currency,
  className = "",
  suffix = " / чел.",
}: TourPriceProps) {
  if (price == null || price <= 0) {
    return null;
  }

  const showOriginal =
    originalPrice != null && originalPrice > 0 && originalPrice > price;

  return (
    <span className={className}>
      {showOriginal ? (
        <>
          <span className="mr-2 text-sm font-normal text-stone-400 line-through">
            {formatPrice(originalPrice, currency)}
          </span>
          <span className="text-red-700">{formatPrice(price, currency)}</span>
        </>
      ) : (
        formatPrice(price, currency)
      )}
      {suffix ? <span className="text-sm font-normal text-stone-500">{suffix}</span> : null}
    </span>
  );
}

export function BurningTourBadge({ compact = false }: { compact?: boolean }) {
  return (
    <span
      className={`rounded-full bg-red-500/95 font-semibold text-white shadow-sm ${
        compact ? "px-2 py-0.5 text-xs" : "px-3 py-1 text-xs"
      }`}
    >
      Горящий тур
    </span>
  );
}
