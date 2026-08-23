import Link from "next/link";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { PERM } from "@/lib/management-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { listManagementBookings, listManagementTours, type BookingStatus } from "@/lib/api";
import { formatDate } from "@/lib/format";

export const dynamic = "force-dynamic";

const STATUSES: { value: "" | BookingStatus; label: string }[] = [
  { value: "", label: "Все" },
  { value: "NEW", label: "Новая" },
  { value: "CONTACTED", label: "Связались" },
  { value: "CONFIRMED", label: "Подтверждена" },
  { value: "COMPLETED", label: "Завершена" },
  { value: "CANCELLED", label: "Отменена" },
];

function asStatus(raw: string | undefined): "" | BookingStatus {
  const allowed = new Set(STATUSES.map((item) => item.value));
  if (!raw || !allowed.has(raw as BookingStatus)) {
    return "";
  }
  return raw as BookingStatus;
}

export default async function ManagementBookingsPage({
  searchParams,
}: {
  searchParams: Promise<{ status?: string; date_from?: string; date_to?: string }>;
}) {
  const access = await canAccessManagementPage([PERM.bookings]);
  if (!access.ok) {
    return <ManagementNoAccess />;
  }

  const params = await searchParams;
  const status = asStatus(params.status);
  const dateFrom = params.date_from?.trim() ?? "";
  const dateTo = params.date_to?.trim() ?? "";
  const [bookings, tours] = await Promise.all([
    listManagementBookings({
      status,
      date_from: dateFrom,
      date_to: dateTo,
      limit: 100,
    }),
    listManagementTours().catch(() => []),
  ]);
  const tourNameById = new Map(tours.map((t) => [t.id, t.title]));
  const exportQuery = new URLSearchParams();
  if (status) exportQuery.set("status", status);
  if (dateFrom) exportQuery.set("date_from", dateFrom);
  if (dateTo) exportQuery.set("date_to", dateTo);
  const exportHref = `/management/bookings/export${exportQuery.toString() ? `?${exportQuery}` : ""}`;

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="font-serif text-3xl">Заявки</h1>
          <p className="text-sm text-muted">Список последних заявок с фильтрами</p>
        </div>
        <a className="btn-ghost" href={exportHref}>
          Скачать CSV
        </a>
      </div>
      <form className="card grid gap-3 p-4 md:grid-cols-4" method="get">
        <label className="text-sm">
          Статус
          <select className="mt-1 w-full rounded-xl border border-[#d7c8b3] bg-white px-3 py-2" name="status" defaultValue={status}>
            {STATUSES.map((item) => (
              <option key={item.value || "all"} value={item.value}>
                {item.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-sm">
          Дата от
          <input className="mt-1 w-full rounded-xl border border-[#d7c8b3] bg-white px-3 py-2" name="date_from" type="date" defaultValue={dateFrom} />
        </label>
        <label className="text-sm">
          Дата до
          <input className="mt-1 w-full rounded-xl border border-[#d7c8b3] bg-white px-3 py-2" name="date_to" type="date" defaultValue={dateTo} />
        </label>
        <div className="flex items-end">
          <button className="btn-primary w-full" type="submit">
            Показать
          </button>
        </div>
      </form>
      <div className="overflow-x-auto rounded-2xl border border-[#e5d9c8] bg-white">
        <table className="w-full min-w-[640px] text-left text-sm">
          <thead className="bg-[#f6efe4] text-muted">
            <tr>
              <th className="px-4 py-3">Тур</th>
              <th className="px-4 py-3">Контакт</th>
              <th className="px-4 py-3">Дата</th>
              <th className="px-4 py-3">Статус</th>
            </tr>
          </thead>
          <tbody>
            {bookings.map((b) => (
              <tr className="border-t border-[#eee4d6]" key={b.id}>
                <td className="px-4 py-3">
                  <Link className="text-accent hover:underline" href={`/management/bookings/${b.id}`}>
                    {tourNameById.get(b.tour_id) ?? "Тур"}
                  </Link>
                </td>
                <td className="px-4 py-3">
                  {b.contact_name}
                  <div className="text-xs text-muted">{b.contact_phone}</div>
                </td>
                <td className="px-4 py-3">{formatDate(b.created_at)}</td>
                <td className="px-4 py-3">{b.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
