import { BookingStatusForm } from "@/components/management/booking-status-form";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge, bookingStatusVariant } from "@/components/management/status-badge";
import { listManagementBookings, listManagementTours } from "@/lib/api/management";
import { formatDateTime, formatManagementBookingStatus, formatPrice } from "@/lib/format";
import { buildTourTitleMap, tourTitle } from "@/lib/tour-title-map";

export default async function ManagementBookingsPage() {
  const [bookings, tours] = await Promise.all([listManagementBookings(), listManagementTours()]);
  const tourNames = buildTourTitleMap(tours);

  return (
    <ManagementPanel title="Заявки" description={`Всего ${bookings.length}. Меняйте статус после связи с клиентом.`}>
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
            <ManagementEmptyRow colSpan={5}>Заявок пока нет.</ManagementEmptyRow>
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
