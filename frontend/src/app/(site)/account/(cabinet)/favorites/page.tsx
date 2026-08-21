import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { EmptyState } from "@/components/empty-state";
import { TourCard } from "@/components/tour-card";
import { getTour } from "@/lib/api/tours";
import { fetchFavoriteTourIds } from "@/lib/auth/user-features";
import { getAuthToken } from "@/lib/auth/session";

export const metadata: Metadata = {
  title: "Избранное",
};

export default async function FavoritesPage() {
  const token = await getAuthToken();
  if (!token) {
    redirect("/account/login?returnUrl=%2Faccount%2Ffavorites");
  }

  let tourIds: string[] = [];
  try {
    tourIds = await fetchFavoriteTourIds();
  } catch {
    redirect("/account/login?returnUrl=%2Faccount%2Ffavorites");
  }

  const tours = (
    await Promise.all(
      tourIds.map(async (id) => {
        try {
          return await getTour(id);
        } catch {
          return null;
        }
      }),
    )
  ).filter(Boolean);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900">Избранное</h1>
        <p className="mt-2 text-sm text-stone-600">
          Туры, которые вы сохранили для сравнения и бронирования
        </p>
      </div>

      {tours.length === 0 ? (
        <EmptyState
          title="Пока нет сохранённых туров"
          description="Нажмите ♡ на карточке тура — он появится здесь, когда будете готовы бронировать."
          actionHref="/search"
          actionLabel="Перейти к поиску"
          secondaryHref="/account/trips"
          secondaryLabel="Мои поездки"
        />
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {tours.map((tour) => (tour ? <TourCard key={tour.id} tour={tour} /> : null))}
        </div>
      )}
    </div>
  );
}
