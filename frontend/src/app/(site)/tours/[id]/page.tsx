import { notFound, permanentRedirect } from "next/navigation";
import { BookingForm } from "@/components/booking-form";
import { Breadcrumbs } from "@/components/breadcrumbs";
import { MobileBookingCTA } from "@/components/mobile-booking-cta";
import { SlotsBadge } from "@/components/slots-badge";
import { TourImage } from "@/components/tour-image";
import { ApiError } from "@/lib/api/client";
import {
  formatDateRange,
  formatPrice,
  formatReviewCount,
  formatTourDuration,
} from "@/lib/format";
import { includedInTour } from "@/lib/site-content";
import { TourViewTracker } from "@/components/tour-view-tracker";
import { getCachedTour, getCachedTourReviews } from "@/lib/api/tour-page";
import { TourStructuredData } from "@/components/structured-data";
import { getSessionUser } from "@/lib/auth/session";
import { toBookingProfile } from "@/lib/auth/user-features";
import { FavoriteButton } from "@/components/favorite-button";
import { CompanyReply } from "@/components/testimonial-card";
import { TourRecommendations } from "@/components/tour-recommendations";
import { isUuidParam, tourPath, tourSeoDescription, tourSeoTitle } from "@/lib/tour-path";
import { isRegularTour, isTourSoldOut } from "@/lib/api/tours";

type TourPageProps = {
  params: Promise<{ id: string }>;
};

async function loadTourPageData(id: string) {
  try {
    const tour = await getCachedTour(id);
    const reviews = await getCachedTourReviews(tour.id);
    return { tour, reviews };
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) {
      notFound();
    }
    throw err;
  }
}

export async function generateMetadata({ params }: TourPageProps) {
  const { id } = await params;

  try {
    const tour = await getCachedTour(id);
    const title = tourSeoTitle(tour);
    const description = tourSeoDescription(tour);
    const images = tour.images?.filter(Boolean).slice(0, 3) ?? [];
    const canonical = tourPath(tour);
    return {
      title,
      description,
      alternates: {
        canonical,
      },
      openGraph: {
        title,
        description,
        type: "website",
        url: canonical,
        ...(images.length > 0 ? { images: images.map((url) => ({ url })) } : {}),
      },
      twitter: {
        card: "summary_large_image",
        title,
        description,
        ...(images[0] ? { images: [images[0]] } : {}),
      },
    };
  } catch {
    return { title: "Тур" };
  }
}

function StarRating({ rating }: { rating: number }) {
  return (
    <span className="inline-flex gap-0.5 text-amber-500" aria-label={`Оценка ${rating} из 5`}>
      {Array.from({ length: 5 }).map((_, index) => (
        <span key={index} className={index < rating ? "opacity-100" : "opacity-25"}>
          ★
        </span>
      ))}
    </span>
  );
}

export default async function TourPage({ params }: TourPageProps) {
  const { id } = await params;
  const [sessionUser, pageData] = await Promise.all([getSessionUser(), loadTourPageData(id)]);
  const { tour, reviews } = pageData;
  const profile = toBookingProfile(sessionUser);
  const soldOut = isTourSoldOut(tour);
  const regular = isRegularTour(tour);
  const duration = regular ? "" : formatTourDuration(tour.date_start, tour.date_end);
  const avgRating =
    reviews.data.length > 0
      ? reviews.data.reduce((sum, r) => sum + r.rating, 0) / reviews.data.length
      : null;

  if (isUuidParam(id) && tour.slug && tour.slug !== id) {
    permanentRedirect(tourPath(tour));
  }

  return (
    <>
      <TourStructuredData tour={tour} />
      <TourViewTracker tourId={tour.id} title={tour.title} />
      <div className="mx-auto max-w-6xl px-4 py-8 pb-28 sm:py-10 lg:pb-10">
        <Breadcrumbs
          items={[
            { name: "Главная", href: "/" },
            { name: "Туры", href: "/search" },
            { name: tour.title },
          ]}
        />

        <div className="grid gap-8 lg:grid-cols-[1.4fr_1fr] lg:items-start">
          <section className="space-y-6">
            <div
              className={`overflow-hidden rounded-2xl border bg-white shadow-sm ${
                tour.is_hot ? "border-amber-200/80 ring-1 ring-amber-100" : "border-stone-200"
              }`}
            >
              <div className="relative">
                <TourImage
                  tour={tour}
                  priority
                  className="aspect-[21/9] w-full sm:aspect-[2/1]"
                  sizes="(max-width: 1024px) 100vw, 60vw"
                />
                <div className="absolute right-4 top-4 z-10">
                  <FavoriteButton tourId={tour.id} />
                </div>
              </div>

              <div className="p-6 sm:p-8">
                <div className="mb-4 flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <h1 className="font-display text-3xl font-semibold tracking-tight sm:text-4xl">
                      {tour.title}
                    </h1>
                    {tour.location ? <p className="mt-2 text-stone-500">{tour.location}</p> : null}
                    {avgRating ? (
                      <p className="mt-2 flex items-center gap-2 text-sm text-stone-600">
                        <StarRating rating={Math.round(avgRating)} />
                        <span>
                          {avgRating.toFixed(1)} · {formatReviewCount(reviews.meta.total)}
                        </span>
                      </p>
                    ) : null}
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {regular ? (
                      <span className="rounded-full bg-stone-100 px-3 py-1 text-xs font-medium text-stone-700">
                        Регулярный тур
                      </span>
                    ) : null}
                    {duration ? (
                      <span className="rounded-full bg-stone-100 px-3 py-1 text-xs font-medium text-stone-700">
                        {duration}
                      </span>
                    ) : null}
                    {tour.is_hot ? (
                      <span className="rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold text-amber-900">
                        Популярный
                      </span>
                    ) : null}
                    <SlotsBadge slotsLeft={tour.slots_left} />
                  </div>
                </div>

                <dl className={`mb-8 grid gap-3 rounded-xl bg-stone-50 p-4 ${regular ? "sm:grid-cols-2" : "sm:grid-cols-4"}`}>
                  {regular ? null : (
                    <>
                      <div>
                        <dt className="text-xs font-medium uppercase tracking-wide text-stone-500">Даты</dt>
                        <dd className="mt-1 text-sm font-medium text-stone-900">
                          {formatDateRange(tour.date_start, tour.date_end)}
                        </dd>
                      </div>
                      <div>
                        <dt className="text-xs font-medium uppercase tracking-wide text-stone-500">Длительность</dt>
                        <dd className="mt-1 text-sm font-medium text-stone-900">{duration || "—"}</dd>
                      </div>
                      <div>
                        <dt className="text-xs font-medium uppercase tracking-wide text-stone-500">Стоимость</dt>
                        <dd className="mt-1 text-sm font-medium text-stone-900">
                          {formatPrice(tour.price, tour.currency)} / чел.
                        </dd>
                      </div>
                    </>
                  )}
                  <div>
                    <dt className="text-xs font-medium uppercase tracking-wide text-stone-500">Группа</dt>
                    <dd className="mt-1 text-sm font-medium text-stone-900">до {tour.slots_total} чел.</dd>
                  </div>
                </dl>

                <div className="mb-8">
                  <h2 className="mb-4 text-lg font-semibold text-stone-900">Программа тура</h2>
                  <div className="space-y-3 text-sm leading-7 text-stone-700">
                    {(tour.description || "Подробное описание тура появится позже.")
                      .split("\n")
                      .filter(Boolean)
                      .map((paragraph, index) => (
                        <p key={index} className={paragraph.startsWith("•") ? "pl-1" : ""}>
                          {paragraph}
                        </p>
                      ))}
                  </div>
                </div>

                <div>
                  <h2 className="mb-3 text-lg font-semibold text-stone-900">Что обычно включено</h2>
                  <ul className="grid gap-2 sm:grid-cols-2">
                    {includedInTour.map((item) => (
                      <li key={item} className="flex items-center gap-2 text-sm text-stone-700">
                        <span className="flex size-5 shrink-0 items-center justify-center rounded-full bg-brand-50 text-xs text-brand-800">
                          ✓
                        </span>
                        {item}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>

            <section className="rounded-2xl border border-stone-200 bg-white p-6 shadow-sm sm:p-8">
              <div className="mb-5 flex items-center justify-between gap-3">
                <h2 className="text-xl font-semibold">Отзывы паломников</h2>
                {reviews.data.length > 0 ? (
                  <span className="text-sm text-stone-500">{formatReviewCount(reviews.meta.total)}</span>
                ) : null}
              </div>

              {reviews.data.length === 0 ? (
                <p className="text-sm text-stone-500">
                  Пока нет опубликованных отзывов — будьте первым после поездки.
                </p>
              ) : (
                <div className="space-y-5">
                  {reviews.data.map((review) => (
                    <article
                      key={review.id}
                      className="rounded-xl bg-stone-50 p-4 first:mt-0"
                    >
                      <div className="mb-2 flex flex-wrap items-center justify-between gap-2">
                        <h3 className="font-medium text-stone-900">{review.client_name}</h3>
                        <StarRating rating={review.rating} />
                      </div>
                      <p className="text-sm leading-6 text-stone-600">{review.text}</p>
                      <CompanyReply text={review.company_reply} />
                    </article>
                  ))}
                </div>
              )}
            </section>

            <TourRecommendations tourId={tour.id} />
          </section>

          <BookingForm tour={tour} profile={profile} />
        </div>
      </div>

      <MobileBookingCTA
        disabled={soldOut}
        price={regular ? undefined : (tour.price ?? undefined)}
        currency={tour.currency}
        dateStart={regular ? undefined : (tour.date_start ?? undefined)}
        dateEnd={regular ? undefined : (tour.date_end ?? undefined)}
      />
    </>
  );
}
