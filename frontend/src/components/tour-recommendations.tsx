import { TourCard } from "@/components/tour-card";
import { getTourRecommendations, type Tour } from "@/lib/api/tours";

export async function TourRecommendations({ tourId }: { tourId: string }) {
  const list = await getTourRecommendations(tourId).catch(() => ({ data: [] as Tour[] }));
  const tours = list.data ?? [];
  if (tours.length === 0) {
    return null;
  }

  return (
    <section className="rounded-2xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
      <h2 className="text-xl font-semibold text-stone-900">Другие опубликованные туры</h2>
      <p className="mt-1 text-sm text-stone-500">Только из каталога, без сгенерированных описаний.</p>
      <div className="mt-5 grid gap-4 sm:grid-cols-2">
        {tours.map((tour) => (
          <TourCard key={tour.id} tour={tour} />
        ))}
      </div>
    </section>
  );
}
