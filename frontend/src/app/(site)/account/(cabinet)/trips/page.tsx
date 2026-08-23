import Link from "next/link";
import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { EmptyState } from "@/components/empty-state";
import { formatBookingStatus, formatDateRange, formatPrice } from "@/lib/format";
import { fetchMyBookings } from "@/lib/api/auth";
import { getAuthToken } from "@/lib/auth/session";
import { getTour } from "@/lib/api/tours";

export const metadata: Metadata = {
  title: "Мои поездки",
};

export default async function MyTripsPage() {
  const token = await getAuthToken();
  if (!token) {
    redirect("/account/login?returnUrl=%2Faccount%2Ftrips");
  }

  let bookings;
  try {
    bookings = await fetchMyBookings(token);
  } catch {
    redirect("/account/login?returnUrl=%2Faccount%2Ftrips");
  }

  const tours = await Promise.all(
    bookings.map(async (booking) => {
      try {
        return await getTour(booking.tour_id);
      } catch {
        return null;
      }
    }),
  );

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900">Мои поездки</h1>
        <p className="mt-2 text-sm text-stone-600">Заявки, оформленные под вашим аккаунтом.</p>
      </div>

      {bookings.length === 0 ? (
        <EmptyState
          title="Пока нет заявок"
          description="Найдите тур и оформите заявку — она появится здесь со статусом обработки."
          actionHref="/search"
          actionLabel="Найти тур"
          secondaryHref="/support/chat"
          secondaryLabel="Спросить в чате"
        />
      ) : (
        <div className="space-y-4">
          {bookings.map((booking, index) => {
            const tour = tours[index];
            return (
              <article
                key={booking.id}
                className="rounded-2xl border border-stone-200 bg-white p-5 shadow-sm"
              >
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <h2 className="text-lg font-semibold text-stone-900">{tour?.title ?? "Тур"}</h2>
                    {tour ? (
                      <p className="mt-1 text-sm text-stone-500">
                        {formatDateRange(tour.date_start, tour.date_end)} · {tour.location}
                      </p>
                    ) : null}
                  </div>
                  <span className="rounded-full bg-brand-50 px-3 py-1 text-xs font-medium text-brand-800">
                    {formatBookingStatus(booking.status)}
                  </span>
                </div>

                <dl className="mt-4 grid gap-2 text-sm text-stone-600 sm:grid-cols-3">
                  <div>
                    <dt className="text-stone-500">Участников</dt>
                    <dd className="font-medium text-stone-900">{booking.people_count}</dd>
                  </div>
                  <div>
                    <dt className="text-stone-500">Сумма</dt>
                    <dd className="font-medium text-stone-900">
                      {formatPrice(booking.total_price, tour?.currency ?? "RUB")}
                    </dd>
                  </div>
                  <div>
                    <dt className="text-stone-500">Заявка</dt>
                    <dd className="font-mono text-xs text-stone-900">{booking.id.slice(0, 8)}…</dd>
                  </div>
                </dl>

                <div className="mt-4 rounded-xl bg-stone-50 p-3 text-sm text-stone-600">
                  <p className="font-medium text-stone-800">Оплата</p>
                  <p className="mt-1">
                    Сумма заявки: {formatPrice(booking.total_price, tour?.currency ?? "RUB")}.
                    Онлайн-оплата на сайте не подключена — порядок оплаты уточняет менеджер.
                  </p>
                </div>

                <div className="mt-4 flex flex-wrap gap-4">
                  {tour ? (
                    <Link
                      href={`/tours/${tour.id}`}
                      className="text-sm font-medium text-brand-800 hover:underline"
                    >
                      Открыть тур
                    </Link>
                  ) : null}
                  <Link href="/account/passengers" className="text-sm font-medium text-stone-600 hover:text-brand-800">
                    Пассажиры
                  </Link>
                  <Link href="/support/chat" className="text-sm font-medium text-stone-600 hover:text-brand-800">
                    Вопрос по заявке
                  </Link>
                </div>
              </article>
            );
          })}
        </div>
      )}
    </div>
  );
}
