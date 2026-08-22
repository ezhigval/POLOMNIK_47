import Link from "next/link";
import { OptimizedImage } from "@/components/optimized-image";
import { popularDestinations } from "@/lib/destinations";

export function PopularDestinations() {
  return (
    <section id="tours" className="scroll-mt-24 space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h2 className="font-display text-2xl font-semibold text-stone-900 sm:text-3xl">
            Популярные направления
          </h2>
          <p className="mt-2 text-sm text-stone-600 sm:text-base">
            Классические маршруты, которые выбирают чаще всего.
          </p>
        </div>
        <Link href="/search" className="btn-secondary text-sm">
          Все туры
        </Link>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {popularDestinations.map((destination) => (
          <Link
            key={destination.id}
            href={`/search?destination=${destination.id}`}
            className="group overflow-hidden rounded-2xl border border-stone-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md"
          >
            <div className="relative aspect-[4/3] overflow-hidden">
              <OptimizedImage
                src={destination.image}
                alt={destination.label}
                fill
                sizes="(max-width: 768px) 100vw, 25vw"
                className="object-cover transition duration-500 group-hover:scale-105"
                fallbackClassName="size-full object-cover transition duration-500 group-hover:scale-105"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-brand-950/80 via-transparent to-transparent" />
              <div className="absolute inset-x-0 bottom-0 p-4 text-white">
                <p className="font-semibold">{destination.label}</p>
                <p className="text-xs text-brand-100">{destination.region}</p>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
