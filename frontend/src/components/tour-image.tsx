import { OptimizedImage } from "@/components/optimized-image";
import type { Tour } from "@/lib/api/tours";
import { getTourCoverUrl } from "@/lib/tour-cover";

type TourImageProps = {
  tour: Tour;
  priority?: boolean;
  className?: string;
  overlay?: boolean;
  sizes?: string;
};

export function TourImage({
  tour,
  priority = false,
  className = "",
  overlay = false,
  sizes = "(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw",
}: TourImageProps) {
  const imageUrl = getTourCoverUrl(tour);

  if (imageUrl) {
    return (
      <div className={`group/img relative overflow-hidden bg-stone-200 ${className}`}>
        <OptimizedImage
          src={imageUrl}
          alt={tour.title}
          fill
          sizes={sizes}
          priority={priority}
          className="object-cover transition duration-500 group-hover/img:scale-105"
          fallbackClassName="size-full object-cover transition duration-500 group-hover/img:scale-105"
        />
        {overlay ? (
          <div className="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent" />
        ) : null}
      </div>
    );
  }

  return (
    <div
      className={`flex items-end overflow-hidden bg-gradient-to-br from-brand-100 via-stone-100 to-amber-50 p-4 ${className}`}
      aria-hidden
    >
      <div>
        <p className="text-xs font-medium uppercase tracking-wider text-brand-700/70">Паломничество</p>
        <p className="line-clamp-2 font-medium text-stone-700">{tour.location || tour.title}</p>
      </div>
    </div>
  );
}
