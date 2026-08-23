import Link from "next/link";
import { SlotsBadge } from "@/components/slots-badge";
import {
  formatDateRange,
  formatPrice,
  formatTourDuration,
  getSlotsAvailability,
} from "@/lib/format";
import type { Tour } from "@/lib/api/tours";

type TourScheduleProps = {
  tours: Tour[];
};

export function TourSchedule({ tours }: TourScheduleProps) {
  return (
    <div className="overflow-hidden rounded-2xl border border-stone-200 bg-white shadow-sm">
      <table className="hidden w-full text-left text-sm lg:table">
        <thead className="bg-stone-50 text-xs font-medium uppercase tracking-wide text-stone-500">
          <tr>
            <th className="px-4 py-3">Даты</th>
            <th className="px-4 py-3">Длительность</th>
            <th className="px-4 py-3">Тур</th>
            <th className="px-4 py-3">Стоимость</th>
            <th className="px-4 py-3">Места</th>
            <th className="px-4 py-3" />
          </tr>
        </thead>
        <tbody>
          {tours.map((tour) => (
            <ScheduleRow key={tour.id} tour={tour} />
          ))}
        </tbody>
      </table>
      <ul className="divide-y divide-stone-100 lg:hidden">
        {tours.map((tour) => (
          <li key={tour.id}>
            <ScheduleCard tour={tour} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function ScheduleRow({ tour }: { tour: Tour }) {
  const soldOut = getSlotsAvailability(tour.slots_left) === "sold_out";
  const duration = formatTourDuration(tour.date_start, tour.date_end);

  return (
    <tr className="border-t border-stone-100 align-top">
      <td className="whitespace-nowrap px-4 py-4 text-stone-800">
        {formatDateRange(tour.date_start, tour.date_end)}
      </td>
      <td className="whitespace-nowrap px-4 py-4 text-stone-600">{duration || "—"}</td>
      <td className="px-4 py-4">
        <Link href={`/tours/${tour.id}`} className="font-medium text-stone-900 hover:text-brand-800">
          {tour.title}
        </Link>
        {tour.location ? <p className="mt-1 text-xs text-stone-500">{tour.location}</p> : null}
      </td>
      <td className="whitespace-nowrap px-4 py-4 font-medium text-stone-900">
        {formatPrice(tour.price, tour.currency)}
        <span className="ml-1 text-xs font-normal text-stone-500">/ чел.</span>
      </td>
      <td className="px-4 py-4">
        <SlotsBadge slotsLeft={tour.slots_left} />
      </td>
      <td className="px-4 py-4 text-right">
        <Link
          href={`/tours/${tour.id}`}
          className={`btn-primary px-4 py-2 ${soldOut ? "pointer-events-none opacity-50" : ""}`}
          aria-disabled={soldOut}
        >
          {soldOut ? "Мест нет" : "Заявка"}
        </Link>
      </td>
    </tr>
  );
}

function ScheduleCard({ tour }: { tour: Tour }) {
  const soldOut = getSlotsAvailability(tour.slots_left) === "sold_out";
  const duration = formatTourDuration(tour.date_start, tour.date_end);

  return (
    <article className="flex flex-col gap-3 p-4">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <div>
          <Link href={`/tours/${tour.id}`} className="font-medium text-stone-900 hover:text-brand-800">
            {tour.title}
          </Link>
          {tour.location ? <p className="mt-1 text-xs text-stone-500">{tour.location}</p> : null}
        </div>
        <SlotsBadge slotsLeft={tour.slots_left} />
      </div>
      <dl className="grid grid-cols-2 gap-2 text-sm text-stone-600">
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-400">Даты</dt>
          <dd className="text-stone-800">{formatDateRange(tour.date_start, tour.date_end)}</dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-400">Длительность</dt>
          <dd className="text-stone-800">{duration || "—"}</dd>
        </div>
        <div className="col-span-2">
          <dt className="text-xs uppercase tracking-wide text-stone-400">Стоимость</dt>
          <dd className="font-medium text-stone-900">
            {formatPrice(tour.price, tour.currency)}
            <span className="ml-1 text-xs font-normal text-stone-500">/ чел.</span>
          </dd>
        </div>
      </dl>
      <Link
        href={`/tours/${tour.id}`}
        className={`btn-primary w-full ${soldOut ? "pointer-events-none opacity-50" : ""}`}
        aria-disabled={soldOut}
      >
        {soldOut ? "Мест нет" : "Смотреть и записаться"}
      </Link>
    </article>
  );
}
