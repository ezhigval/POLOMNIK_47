import Link from "next/link";
import { Suspense } from "react";
import type { Metadata } from "next";
import { TripSearchConstructor } from "@/components/trip-search-constructor";
import { ToursSection } from "@/components/tours-section";
import { findDestination } from "@/lib/destinations";
import { parseTourFilters } from "@/lib/tour-filters";

export const metadata: Metadata = {
  title: "Поиск туров",
  description: "Найдите паломнический тур по направлению, датам и количеству участников.",
  alternates: { canonical: "/search" },
};

type SearchPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function SearchPage({ searchParams }: SearchPageProps) {
  const params = await searchParams;
  const filters = parseTourFilters(params);
  const destination = filters.destination ? findDestination(filters.destination) : undefined;

  return (
    <div className="mx-auto max-w-6xl space-y-8 px-4 py-8 sm:py-10">
      <div className="space-y-3">
        <p className="text-sm font-medium uppercase tracking-widest text-brand-800">Поиск</p>
        <h1 className="font-display text-3xl font-semibold text-stone-900 sm:text-4xl">
          {destination ? `Туры: ${destination.label}` : "Подбор паломнического тура"}
        </h1>
        <p className="max-w-2xl text-sm text-stone-600 sm:text-base">
          Укажите направление, даты и количество человек — покажем доступные поездки с учётом
          свободных мест.
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
    </div>
  );
}
