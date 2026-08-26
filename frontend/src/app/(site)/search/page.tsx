import Link from "next/link";
import { Suspense } from "react";
import type { Metadata } from "next";
import { CatalogStickyCTA } from "@/components/catalog-sticky-cta";
import { TripSearchConstructor } from "@/components/trip-search-constructor";
import { ToursSection } from "@/components/tours-section";
import { Breadcrumbs } from "@/components/breadcrumbs";
import { hasActiveFilters, parseTourFilters } from "@/lib/tour-filters";
import { buildPublicPageMetadata } from "@/lib/seo-metadata";
import { siteConfig } from "@/lib/site-config";

type SearchPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export async function generateMetadata({ searchParams }: SearchPageProps): Promise<Metadata> {
  const params = await searchParams;
  const filters = parseTourFilters(params);
  const filtered = hasActiveFilters(filters);
  const description = `Расписание паломнических туров из ${siteConfig.departureCity}: даты, стоимость и длительность.`;

  return {
    ...buildPublicPageMetadata({
      title: "Расписание туров",
      description,
      canonical: "/search",
    }),
    robots: filtered ? { index: false, follow: true } : { index: true, follow: true },
  };
}

export default async function SearchPage({ searchParams }: SearchPageProps) {
  const params = await searchParams;
  const filters = parseTourFilters(params);

  return (
    <div className="mx-auto max-w-6xl space-y-8 px-4 py-8 pb-24 sm:py-10 lg:pb-10">
      <Breadcrumbs
        items={[
          { name: "Главная", href: "/" },
          { name: "Расписание туров" },
        ]}
      />
      <div className="space-y-3">
        <p className="text-sm font-medium uppercase tracking-widest text-brand-800">Туры</p>
        <h1 className="font-display text-3xl font-semibold text-stone-900 sm:text-4xl">
          Расписание паломнических туров
        </h1>
        <p className="max-w-2xl text-sm text-stone-600 sm:text-base">
          Даты, стоимость и длительность — из карточки тура. Выезд из {siteConfig.departureCity}.
          Фильтр по датам и числу мест.
        </p>
      </div>

      <TripSearchConstructor initialValues={filters} />

      <Suspense fallback={<div className="h-64 animate-pulse rounded-2xl bg-stone-200" />}>
        <ToursSection searchParams={params} basePath="/search" showPopular={false} showFilters={false} />
      </Suspense>

      <div className="rounded-2xl border border-stone-200 bg-stone-50 p-5 text-sm text-stone-600">
        Не нашли подходящий тур?{" "}
        <Link href="/support/chat" className="font-medium text-brand-800 underline-offset-2 hover:underline">
          Напишите в чат поддержки
        </Link>{" "}
        или посмотрите{" "}
        <Link href="/support" className="font-medium text-brand-800 underline-offset-2 hover:underline">
          справочник
        </Link>
        .
      </div>
      <CatalogStickyCTA />
    </div>
  );
}
