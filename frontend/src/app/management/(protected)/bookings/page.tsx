import { BookingStatusForm } from "@/components/management/booking-status-form";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge, bookingStatusVariant } from "@/components/management/status-badge";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { PERM } from "@/lib/management-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { listManagementBookings, listManagementTours } from "@/lib/api/management";
import { formatDateTime, formatManagementBookingStatus, formatPrice } from "@/lib/format";
import { buildTourTitleMap, tourTitle } from "@/lib/tour-title-map";

const BOOKING_STATUSES = ["NEW", "CONTACTED", "CONFIRMED", "COMPLETED", "CANCELLED"] as const;

type PageProps = {
  searchParams: Promise<{ status?: string; date_from?: string; date_to?: string }>;
};

export default async function ManagementBookingsPage({ searchParams }: PageProps) {
  if (!(await canAccessManagementPage([PERM.bookings]))) {
    return <ManagementNoAccess />;
  }

  const filters = await searchParams;
  const status = BOOKING_STATUSES.includes(filters.status as (typeof BOOKING_STATUSES)[number])
    ? filters.status
    : undefined;
  const dateFrom = filters.date_from?.trim() || undefined;
  const dateTo = filters.date_to?.trim() || undefined;

  const [bookings, tours] = await Promise.all([
    listManagementBookings({
      status,
      date_from: dateFrom,
      date_to: dateTo,
      limit: 100,
    }),
    listManagementTours().catch(() => []),
  ]);
  const tourNames = buildTourTitleMap(tours);
  const exportQuery = new URLSearchParams();
  if (status) exportQuery.set("status", status);
  if (dateFrom) exportQuery.set("date_from", dateFrom);
  if (dateTo) exportQuery.set("date_to", dateTo);
  const exportHref = `/management/bookings/export${exportQuery.toString() ? `?${exportQuery}` : ""}`;

  return (
    <ManagementPanel
      title="Заявки"
      description={`Показано ${bookings.length}. Меняйте статус после связи с клиентом. CSV — те же фильтры, все совпадения, не только эта страница.`}
    >
      <form method="get" className="flex flex-wrap items-end gap-3 border-b border-stone-100 px-4 py-4">
        <label className="text-sm">
          <span className="mb-1 block text-stone-500">Статус</span>
          <select name="status" defaultValue={status ?? ""} className="input-field min-w-[180px]">
            <option value="">Все</option>
            {BOOKING_STATUSES.map((value) => (
              <option key={value} value={value}>
                {formatManagementBookingStatus(value)}
              </option>
            ))}
          </select>
        </label>
        <label className="text-sm">
          <span className="mb-1 block text-stone-500">С даты</span>
          <input type="date" name="date_from" defaultValue={dateFrom ?? ""} className="input-field" />
        </label>
        <label className="text-sm">
          <span className="mb-1 block text-stone-500">По дату</span>
          <input type="date" name="date_to" defaultValue={dateTo ?? ""} className="input-field" />
        </label>
        <button type="submit" className="btn-secondary">
          Фильтр
        </button>
        <a href={exportHref} className="btn-primary">
          Скачать CSV
        </a>
      </form>
      <ManagementTable>
        <ManagementTableHead>
          <ManagementTh>Клиент</ManagementTh>
          <ManagementTh>Тур</ManagementTh>
          <ManagementTh>Детали</ManagementTh>
          <ManagementTh>Сумма</ManagementTh>
          <ManagementTh>Статус</ManagementTh>
        </ManagementTableHead>
        <tbody>
          {bookings.length === 0 ? (
            <ManagementEmptyRow colSpan={5}>Заявок по фильтру нет.</ManagementEmptyRow>
          ) : (
            bookings.map((booking) => (
              <tr key={booking.id} className="border-b border-stone-100 align-top last:border-0">
                <td className="px-4 py-4">
                  <div className="font-medium text-stone-900">{booking.name}</div>
                  <div className="text-stone-500">{booking.phone}</div>
                  {booking.email ? <div className="text-stone-500">{booking.email}</div> : null}
                  <div className="mt-1 font-mono text-xs text-stone-400">{booking.id}</div>
                </td>
                <td className="px-4 py-4">
                  <div className="font-medium">{tourTitle(tourNames, booking.tour_id)}</div>
                  <div className="font-mono text-xs text-stone-400">{booking.tour_id}</div>
                </td>
                <td className="px-4 py-4">
                  <div>{booking.people_count} чел.</div>
                  <div className="text-stone-500">{formatDateTime(booking.created_at)}</div>
                  {booking.comment ? (
                    <p className="mt-2 max-w-xs rounded-lg bg-stone-50 p-2 text-stone-600">{booking.comment}</p>
                  ) : null}
                  {booking.overbooked ? (
                    <StatusBadge variant="warning">Overbooking</StatusBadge>
                  ) : null}
                </td>
                <td className="px-4 py-4 font-medium">{formatPrice(booking.total_price, "₽")}</td>
                <td className="px-4 py-4">
                  <div className="mb-2">
                    <StatusBadge variant={bookingStatusVariant(booking.status)}>
                      {formatManagementBookingStatus(booking.status)}
                    </StatusBadge>
                  </div>
                  <BookingStatusForm bookingId={booking.id} currentStatus={booking.status} />
                </td>
              </tr>
            ))
          )}
        </tbody>
      </ManagementTable>
    </ManagementPanel>
  );
}
