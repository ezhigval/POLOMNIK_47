import Link from "next/link";
import { SlotsBadge } from "@/components/slots-badge";
import {
  formatDateRange,
  formatPrice,
  formatTourDuration,
} from "@/lib/format";
import { isRegularTour, isTourSoldOut, tourShowsPrice, type Tour } from "@/lib/api/tours";
import { tourPath } from "@/lib/tour-path";

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
  const soldOut = isTourSoldOut(tour);
  const regular = isRegularTour(tour);
  const showPrice = tourShowsPrice(tour);
  const duration = regular ? "" : formatTourDuration(tour.date_start, tour.date_end);

  return (
    <tr className="border-t border-stone-100 align-top">
      <td className="whitespace-nowrap px-4 py-4 text-stone-800">
        {regular ? "Регулярный тур" : formatDateRange(tour.date_start, tour.date_end)}
      </td>
      <td className="whitespace-nowrap px-4 py-4 text-stone-600">{duration || "—"}</td>
      <td className="px-4 py-4">
        <Link href={tourPath(tour)} className="font-medium text-stone-900 hover:text-brand-800">
          {tour.title}
        </Link>
        {tour.is_hot ? (
          <span className="ml-2 inline-flex rounded-full bg-amber-100 px-2 py-0.5 text-xs font-semibold text-amber-900">
            Популярный
          </span>
        ) : null}
        {tour.location ? <p className="mt-1 text-xs text-stone-500">{tour.location}</p> : null}
      </td>
      <td className="whitespace-nowrap px-4 py-4 font-medium text-stone-900">
        {showPrice ? (
          <>
            {formatPrice(tour.price, tour.currency)}
            <span className="ml-1 text-xs font-normal text-stone-500">/ чел.</span>
          </>
        ) : (
          "—"
        )}
      </td>
      <td className="px-4 py-4">
        <SlotsBadge slotsLeft={tour.slots_left} />
      </td>
      <td className="px-4 py-4 text-right">
        <Link
          href={tourPath(tour)}
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
  const soldOut = isTourSoldOut(tour);
  const regular = isRegularTour(tour);
  const showPrice = tourShowsPrice(tour);
  const duration = regular ? "" : formatTourDuration(tour.date_start, tour.date_end);

  return (
    <article className="flex flex-col gap-3 p-4">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <div>
          <Link href={tourPath(tour)} className="font-medium text-stone-900 hover:text-brand-800">
            {tour.title}
          </Link>
          {tour.is_hot ? (
            <span className="mt-1 inline-flex rounded-full bg-amber-100 px-2 py-0.5 text-xs font-semibold text-amber-900">
              Популярный
            </span>
          ) : null}
          {tour.location ? <p className="mt-1 text-xs text-stone-500">{tour.location}</p> : null}
        </div>
        <SlotsBadge slotsLeft={tour.slots_left} />
      </div>
      {regular && !showPrice ? (
        <p className="text-sm text-stone-600">Регулярный тур</p>
      ) : (
        <dl className="grid grid-cols-2 gap-2 text-sm text-stone-600">
          {regular ? null : (
            <>
              <div>
                <dt className="text-xs uppercase tracking-wide text-stone-400">Даты</dt>
                <dd className="text-stone-800">{formatDateRange(tour.date_start, tour.date_end)}</dd>
              </div>
              <div>
                <dt className="text-xs uppercase tracking-wide text-stone-400">Длительность</dt>
                <dd className="text-stone-800">{duration || "—"}</dd>
              </div>
            </>
          )}
          {showPrice ? (
            <div className="col-span-2">
              <dt className="text-xs uppercase tracking-wide text-stone-400">Стоимость</dt>
              <dd className="font-medium text-stone-900">
                {formatPrice(tour.price, tour.currency)}
                <span className="ml-1 text-xs font-normal text-stone-500">/ чел.</span>
              </dd>
            </div>
          ) : regular ? (
            <div className="col-span-2">
              <p className="text-sm text-stone-600">Регулярный тур</p>
            </div>
          ) : null}
        </dl>
      )}
      <Link
        href={tourPath(tour)}
        className={`btn-primary w-full ${soldOut ? "pointer-events-none opacity-50" : ""}`}
        aria-disabled={soldOut}
      >
        {soldOut ? "Мест нет" : "Смотреть и записаться"}
      </Link>
    </article>
  );
}
