import Link from "next/link";
import {
  listManagementBookings,
  listManagementCmsPagesOrEmpty,
  listManagementIntegrationReferences,
  listManagementNews,
  listManagementOutboxEvents,
  listManagementReviews,
  listManagementSupportThreads,
  listManagementTours,
  getManagementSession,
} from "@/lib/api/management";
import { PERM, sessionHasPermission } from "@/lib/management-access";

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

function settledList<T>(result: PromiseSettledResult<T[]>, fallback: T[] = []): T[] {
  return result.status === "fulfilled" ? result.value : fallback;
}

export default async function ManagementDashboardPage() {
  const session = await getManagementSession().catch(() => ({
    full_admin: false,
    permissions: [] as string[],
  }));
  const can = (perm: string) => sessionHasPermission(session, perm);

  const [
    toursResult,
    bookingsResult,
    reviewsResult,
    supportResult,
    integrationRefsResult,
    outboxEventsResult,
    cmsResult,
    newsResult,
  ] = await Promise.allSettled([
    can(PERM.tours) ? listManagementTours() : Promise.resolve([]),
    can(PERM.bookings) ? listManagementBookings() : Promise.resolve([]),
    can(PERM.content) ? listManagementReviews() : Promise.resolve([]),
    can(PERM.support) ? listManagementSupportThreads() : Promise.resolve([]),
    can(PERM.integrations) ? listManagementIntegrationReferences() : Promise.resolve([]),
    can(PERM.integrations) ? listManagementOutboxEvents({ status: "pending" }) : Promise.resolve([]),
    can(PERM.content) ? listManagementCmsPagesOrEmpty() : Promise.resolve({ pages: [], unavailable: true }),
    can(PERM.content) ? listManagementNews() : Promise.resolve([]),
  ]);

  const tours = settledList(toursResult);
  const bookings = settledList(bookingsResult);
  const reviews = settledList(reviewsResult);
  const supportThreads = settledList(supportResult);
  const integrationRefs = settledList(integrationRefsResult);
  const outboxEvents = settledList(outboxEventsResult);
  const news = settledList(newsResult);
  const cms =
    cmsResult.status === "fulfilled" ? cmsResult.value : { pages: [], unavailable: true };
  const cmsPages = cms.pages;

  const newBookings = bookings.filter((booking) => booking.status === "NEW").length;
  const pendingReviews = reviews.filter((review) => !review.is_approved).length;
  const openSupport = supportThreads.filter((thread) => thread.status === "open").length;
  const failedSync = integrationRefs.filter((ref) => ref.sync_status === "failed").length;
  const pendingOutbox = outboxEvents.length;

  const cards: StatCardProps[] = [];
  if (can(PERM.content)) {
    cards.push(
      {
        title: "Главная",
        value: cmsPages.filter((page) => page.slug === "home").length,
        hint: cms.unavailable
          ? "недоступно"
          : cmsPages.some((page) => page.slug === "home")
            ? "блоки и SEO"
            : "ещё не создана",
        href: "/management/content",
        accent: cms.unavailable || !cmsPages.some((page) => page.slug === "home"),
      },
      {
        title: "Новости",
        value: news.length,
        hint: `${news.filter((item) => item.is_published).length} опубликовано`,
        href: "/management/news",
      },
    );
  }
  if (can(PERM.tours)) {
    cards.push({
      title: "Туры",
      value: tours.length,
      hint: "активных и скрытых",
      href: "/management/tours",
    });
  }
  if (can(PERM.bookings)) {
    cards.push({
      title: "Новые заявки",
      value: newBookings,
      hint: `из ${bookings.length} всего`,
      href: "/management/bookings",
      accent: newBookings > 0,
    });
  }
  if (can(PERM.support)) {
    cards.push({
      title: "Поддержка",
      value: openSupport,
      hint: `из ${supportThreads.length} всего`,
      href: "/management/support",
      accent: openSupport > 0,
    });
  }
  if (can(PERM.content)) {
    cards.push({
      title: "Отзывы на модерации",
      value: pendingReviews,
      hint: `из ${reviews.length} всего`,
      href: "/management/reviews",
      accent: pendingReviews > 0,
    });
  }
  if (can(PERM.integrations) || can(PERM.stats)) {
    cards.push({
      title: "Синхронизация",
      value: integrationRefs.length,
      hint:
        failedSync > 0
          ? `${failedSync} с ошибкой`
          : pendingOutbox > 0
            ? `${pendingOutbox} в очереди`
            : "Bitrix / 1С",
      href: "/management/integrations",
      accent: failedSync > 0 || pendingOutbox > 0,
    });
  }

  const actions: Array<{ href: string; label: string; primary?: boolean }> = [];
  if (can(PERM.content)) {
    actions.push(
      { href: "/management/content", label: "Редактировать главную", primary: true },
      { href: "/management/news", label: "Новости" },
      { href: "/management/smm", label: "Контент-план" },
    );
  }
  if (can(PERM.tours)) {
    actions.push({ href: "/management/tours", label: "Создать тур" });
  }
  if (can(PERM.bookings)) {
    actions.push({ href: "/management/bookings", label: "Обработать заявки" });
  }
  if (can(PERM.support)) {
    actions.push({ href: "/management/support", label: "Поддержка" });
  }
  if (can(PERM.stats)) {
    actions.push({ href: "/management/ai", label: "Дайджест и watchdog" });
  }
  if (can(PERM.content)) {
    actions.push({ href: "/management/reviews", label: "Модерировать отзывы" });
  }
  if (can(PERM.settingsSite) || can(PERM.recipients) || can(PERM.roles)) {
    actions.push({ href: "/management/settings", label: "Настройки" });
  }

  return (
    <div className="space-y-8">
      {can(PERM.stats) ? (
        <section className="rounded-2xl border border-stone-200 bg-white p-5">
          <h2 className="text-lg font-semibold text-stone-900">Визиты сайта</h2>
          {process.env.NEXT_PUBLIC_YM_ID?.trim() ? (
            <p className="mt-2 text-sm leading-6 text-stone-600">
              Счётчик Яндекс.Метрики подключён (ID в env). Отчёты по визитам — в{" "}
              <a
                href="https://metrika.yandex.ru/"
                target="_blank"
                rel="noreferrer"
                className="font-medium text-brand-800 underline-offset-2 hover:underline"
              >
                кабинете Метрики
              </a>
              . Цифры визитов на этой странице не дублируются.
            </p>
          ) : (
            <p className="mt-2 text-sm leading-6 text-stone-600">
              Подключите Метрику: задайте <code className="rounded bg-stone-100 px-1">NEXT_PUBLIC_YM_ID</code>{" "}
              в prod env и пересоберите фронт. Визиты здесь не выдумываются.
            </p>
          )}
        </section>
      ) : null}

      {cards.length > 0 ? (
        <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {cards.map((card) => (
            <StatCard key={card.href} {...card} />
          ))}
        </section>
      ) : (
        <p className="rounded-2xl border border-stone-200 bg-white p-5 text-sm text-stone-600">
          Для этой роли нет карточек разделов. Если нужны другие права — полный админ выдаёт их в
          «Настройках».
        </p>
      )}

      {actions.length > 0 ? (
        <section className="rounded-2xl border border-stone-200 bg-white p-5">
          <h2 className="mb-4 text-lg font-semibold">Быстрые действия</h2>
          <div className="flex flex-wrap gap-3">
            {actions.map((action) => (
              <Link
                key={`${action.href}-${action.label}`}
                href={action.href}
                className={action.primary ? "btn-primary" : "btn-secondary"}
              >
                {action.label}
              </Link>
            ))}
          </div>
        </section>
      ) : null}
    </div>
  );
}
