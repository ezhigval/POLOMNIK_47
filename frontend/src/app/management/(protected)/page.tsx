import Link from "next/link";
import {
  listManagementBookings,
  listManagementCmsPages,
  listManagementIntegrationReferences,
  listManagementOutboxEvents,
  listManagementReviews,
  listManagementTours,
} from "@/lib/api/management";

type StatCardProps = {
  title: string;
  value: number;
  hint: string;
  href: string;
  accent?: boolean;
};

function StatCard({ title, value, hint, href, accent }: StatCardProps) {
  return (
    <Link
      href={href}
      className={`group rounded-2xl border bg-white p-5 transition hover:shadow-sm ${
        accent ? "border-brand-200 ring-1 ring-brand-100" : "border-stone-200 hover:border-brand-200"
      }`}
    >
      <p className="text-sm text-stone-500">{title}</p>
      <p className="mt-2 text-3xl font-semibold tracking-tight text-stone-900 group-hover:text-brand-800">
        {value}
      </p>
      <p className="mt-1 text-sm text-stone-600">{hint}</p>
    </Link>
  );
}

export default async function ManagementDashboardPage() {
  const [tours, bookings, reviews, integrationRefs, outboxEvents, cmsPages] = await Promise.all([
    listManagementTours(),
    listManagementBookings(),
    listManagementReviews(),
    listManagementIntegrationReferences(),
    listManagementOutboxEvents({ status: "pending" }),
    listManagementCmsPages(),
  ]);

  const newBookings = bookings.filter((booking) => booking.status === "NEW").length;
  const pendingReviews = reviews.filter((review) => !review.is_approved).length;
  const failedSync = integrationRefs.filter((ref) => ref.sync_status === "failed").length;
  const pendingOutbox = outboxEvents.length;

  const cards: StatCardProps[] = [
    {
      title: "Страницы CMS",
      value: cmsPages.length,
      hint: cmsPages.some((page) => page.slug === "home") ? "главная настроена" : "главная не создана",
      href: "/management/content",
      accent: !cmsPages.some((page) => page.slug === "home"),
    },
    {
      title: "Туры",
      value: tours.length,
      hint: "активных и скрытых",
      href: "/management/tours",
    },
    {
      title: "Новые заявки",
      value: newBookings,
      hint: `из ${bookings.length} всего`,
      href: "/management/bookings",
      accent: newBookings > 0,
    },
    {
      title: "Отзывы на модерации",
      value: pendingReviews,
      hint: `из ${reviews.length} всего`,
      href: "/management/reviews",
      accent: pendingReviews > 0,
    },
    {
      title: "Sync записи",
      value: integrationRefs.length,
      hint:
        failedSync > 0
          ? `${failedSync} с ошибкой`
          : pendingOutbox > 0
            ? `${pendingOutbox} в outbox`
            : "Bitrix / 1С",
      href: "/management/integrations",
      accent: failedSync > 0 || pendingOutbox > 0,
    },
  ];

  return (
    <div className="space-y-8">
      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {cards.map((card) => (
          <StatCard key={card.href} {...card} />
        ))}
      </section>

      <section className="rounded-2xl border border-stone-200 bg-white p-5">
        <h2 className="mb-4 text-lg font-semibold">Быстрые действия</h2>
        <div className="flex flex-wrap gap-3">
          <Link href="/management/content" className="btn-primary">
            Редактировать контент
          </Link>
          <Link href="/management/tours" className="btn-secondary">
            Создать тур
          </Link>
          <Link href="/management/bookings" className="btn-secondary">
            Обработать заявки
          </Link>
          <Link href="/management/reviews" className="btn-secondary">
            Модерировать отзывы
          </Link>
        </div>
      </section>
    </div>
  );
}
