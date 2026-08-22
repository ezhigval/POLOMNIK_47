"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { FormEvent, useState } from "react";
import type { TourFilterValues } from "@/lib/tour-filters";

type TourFiltersProps = {
  initial: TourFilterValues;
  basePath?: string;
};

const filterFields = [
  { key: "date_from" as const, label: "Дата с", type: "date" },
  { key: "date_to" as const, label: "Дата по", type: "date" },
  { key: "location" as const, label: "Направление", type: "text", placeholder: "Дивеево" },
  { key: "price_min" as const, label: "Цена от", type: "number" },
  { key: "price_max" as const, label: "Цена до", type: "number" },
];

export function TourFilters({ initial, basePath = "/search" }: TourFiltersProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [expanded, setExpanded] = useState(false);

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const params = new URLSearchParams(searchParams.toString());

    for (const field of ["date_from", "date_to", "price_min", "price_max", "location", "is_hot"] as const) {
      const value = String(formData.get(field) ?? "").trim();
      if (value) {
        params.set(field, value);
      } else {
        params.delete(field);
      }
    }

    params.delete("page");
    const query = params.toString();
    router.push(query ? `${basePath}?${query}` : basePath);
  }

  function onReset() {
    router.push(basePath);
  }

  return (
    <div className="rounded-2xl border border-stone-200 bg-white shadow-sm">
      <div className="flex items-center justify-between gap-3 border-b border-stone-100 px-5 py-4 lg:hidden">
        <p className="font-medium text-stone-900">Фильтры</p>
        <button
          type="button"
          onClick={() => setExpanded((value) => !value)}
          className="btn-secondary px-4 py-1.5 text-xs"
          aria-expanded={expanded}
        >
          {expanded ? "Скрыть" : "Показать"}
        </button>
      </div>

      <form
        onSubmit={onSubmit}
        className={`grid gap-4 p-5 md:grid-cols-2 lg:grid-cols-3 ${expanded ? "grid" : "hidden lg:grid"}`}
        aria-label="Фильтры туров"
      >
        {filterFields.map((field) => (
          <label key={field.key} className="block text-sm">
            <span className="mb-1.5 block font-medium text-stone-700">{field.label}</span>
            <input
              type={field.type}
              name={field.key}
              min={field.type === "number" ? 0 : undefined}
              defaultValue={initial[field.key]}
              placeholder={field.placeholder}
              className="input-field"
            />
          </label>
        ))}

        <label className="flex items-end gap-2.5 pb-2.5 text-sm">
          <input
            type="checkbox"
            name="is_hot"
            value="true"
            defaultChecked={initial.is_hot === "true"}
            className="size-4 rounded border-stone-300 text-brand-800 focus:ring-brand-700"
          />
          <span className="font-medium text-stone-700">Только популярные</span>
        </label>

        <div className="flex flex-wrap gap-3 md:col-span-2 lg:col-span-3">
          <button type="submit" className="btn-primary">
            Применить
          </button>
          <button type="button" onClick={onReset} className="btn-secondary">
            Сбросить
          </button>
        </div>
      </form>
    </div>
  );
}

type ActiveFilterChipsProps = {
  filters: TourFilterValues;
  basePath?: string;
};

const chipLabels: Partial<Record<keyof TourFilterValues, string>> = {
  date_from: "с",
  date_to: "по",
  location: "",
  price_min: "от",
  price_max: "до",
  is_hot: "Популярные",
  q: "Направление",
  people: "Паломников",
  page: "",
};

export function ActiveFilterChips({ filters, basePath = "/search" }: ActiveFilterChipsProps) {
  const chips: { key: string; label: string; removeKey: string }[] = [];

  if (filters.date_from) chips.push({ key: "date_from", label: `Дата ${chipLabels.date_from} ${filters.date_from}`, removeKey: "date_from" });
  if (filters.date_to) chips.push({ key: "date_to", label: `Дата ${chipLabels.date_to} ${filters.date_to}`, removeKey: "date_to" });
  if (filters.location) chips.push({ key: "location", label: filters.location, removeKey: "location" });
  if (filters.price_min) chips.push({ key: "price_min", label: `Цена от ${filters.price_min}`, removeKey: "price_min" });
  if (filters.price_max) chips.push({ key: "price_max", label: `Цена до ${filters.price_max}`, removeKey: "price_max" });
  if (filters.is_hot === "true") chips.push({ key: "is_hot", label: chipLabels.is_hot ?? "Популярные", removeKey: "is_hot" });
  if (filters.q) chips.push({ key: "q", label: `${chipLabels.q}: ${filters.q}`, removeKey: "q" });
  if (filters.people) chips.push({ key: "people", label: `Паломников: ${filters.people}`, removeKey: "people" });

  if (chips.length === 0) return null;

  return (
    <div className="flex flex-wrap items-center gap-2">
      <span className="text-sm text-stone-500">Активные фильтры:</span>
      {chips.map((chip) => (
        <Link
          key={chip.key}
          href={buildRemoveFilterHref(filters, chip.removeKey, basePath)}
          className="inline-flex items-center gap-1 rounded-full bg-brand-50 px-3 py-1 text-xs font-medium text-brand-800 ring-1 ring-brand-100 transition hover:bg-brand-100"
        >
          {chip.label}
          <span aria-hidden>×</span>
        </Link>
      ))}
      <Link href={basePath} className="text-xs text-stone-500 underline-offset-2 hover:underline">
        Сбросить все
      </Link>
    </div>
  );
}

function buildRemoveFilterHref(filters: TourFilterValues, removeKey: string, basePath: string): string {
  const params = new URLSearchParams();
  const skip = new Set(["page", "from", "destination", "min_slots", removeKey]);
  for (const [key, value] of Object.entries(filters)) {
    if (skip.has(key) || !value) continue;
    params.set(key, value);
  }
  const query = params.toString();
  return query ? `${basePath}?${query}` : basePath;
}
