"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { departureCities, popularDestinations } from "@/lib/destinations";
import { siteConfig } from "@/lib/site-config";
import { toSearchParams, type TripSearchValues } from "@/lib/tour-filters";

type TripSearchConstructorProps = {
  initialValues?: TripSearchValues;
  compact?: boolean;
  className?: string;
};

export function TripSearchConstructor({
  initialValues,
  compact = false,
  className = "",
}: TripSearchConstructorProps) {
  const router = useRouter();
  const [from, setFrom] = useState(initialValues?.from ?? siteConfig.departureCity);
  const [destination, setDestination] = useState(initialValues?.destination ?? "");
  const [dateFrom, setDateFrom] = useState(initialValues?.date_from ?? "");
  const [dateTo, setDateTo] = useState(initialValues?.date_to ?? "");
  const [people, setPeople] = useState(initialValues?.people ?? "2");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const params = toSearchParams({
      from,
      destination,
      date_from: dateFrom,
      date_to: dateTo,
      people,
    });
    router.push(`/search?${params.toString()}`);
  }

  return (
    <form
      onSubmit={handleSubmit}
      className={`search-constructor ${compact ? "search-constructor-compact" : ""} ${className}`}
    >
      <div className="search-constructor-grid">
        <label className="search-field">
          <span className="search-field-label">Откуда</span>
          <select
            className="search-field-input"
            value={from}
            onChange={(event) => setFrom(event.target.value)}
          >
            {departureCities.map((city) => (
              <option key={city} value={city}>
                {city}
              </option>
            ))}
          </select>
        </label>

        <label className="search-field">
          <span className="search-field-label">Куда</span>
          <select
            className="search-field-input"
            value={destination}
            onChange={(event) => setDestination(event.target.value)}
          >
            <option value="">Любое направление</option>
            {popularDestinations.map((item) => (
              <option key={item.id} value={item.id}>
                {item.label}
              </option>
            ))}
          </select>
        </label>

        <label className="search-field">
          <span className="search-field-label">Дата с</span>
          <input
            type="date"
            className="search-field-input"
            value={dateFrom}
            onChange={(event) => setDateFrom(event.target.value)}
          />
        </label>

        <label className="search-field">
          <span className="search-field-label">Дата по</span>
          <input
            type="date"
            className="search-field-input"
            value={dateTo}
            onChange={(event) => setDateTo(event.target.value)}
          />
        </label>

        <label className="search-field">
          <span className="search-field-label">Паломников</span>
          <select
            className="search-field-input"
            value={people}
            onChange={(event) => setPeople(event.target.value)}
          >
            {Array.from({ length: 8 }, (_, index) => index + 1).map((count) => (
              <option key={count} value={String(count)}>
                {count} {peopleLabel(count)}
              </option>
            ))}
          </select>
        </label>

        <div className="search-field search-field-action">
          <span className="search-field-label opacity-0">Поиск</span>
          <button type="submit" className="btn-primary w-full py-3 text-base">
            Найти туры
          </button>
        </div>
      </div>
    </form>
  );
}

function peopleLabel(count: number): string {
  if (count === 1) return "человек";
  if (count >= 2 && count <= 4) return "человека";
  return "человек";
}
