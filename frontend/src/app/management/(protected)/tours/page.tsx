import { CreateTourForm } from "@/components/management/create-tour-form";
import { EditTourForm } from "@/components/management/edit-tour-form";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { deleteTourAction } from "@/app/management/actions";
import { listManagementTours } from "@/lib/api/management";
import { formatDateRange, formatPrice, formatTourDuration } from "@/lib/format";
import { tourShowsPrice } from "@/lib/api/tours";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

export default async function ManagementToursPage() {
  if (!(await canAccessManagementPage([PERM.tours]))) {
    return <ManagementNoAccess />;
  }
  const tours = await listManagementTours();

  return (
    <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
      <ManagementPanel title="Туры" description={`${tours.length} в каталоге`}>
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Название</ManagementTh>
            <ManagementTh>Даты</ManagementTh>
            <ManagementTh>Длительность</ManagementTh>
            <ManagementTh>Места</ManagementTh>
            <ManagementTh>Цена</ManagementTh>
            <ManagementTh>Статус</ManagementTh>
            <ManagementTh />
          </ManagementTableHead>
          <tbody>
            {tours.length === 0 ? (
              <ManagementEmptyRow colSpan={7}>Туров пока нет.</ManagementEmptyRow>
            ) : (
              tours.map((tour) => (
                <tr key={tour.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4">
                    <div className="font-medium text-stone-900">{tour.title}</div>
                    <div className="text-stone-500">{tour.location}</div>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap">
                    {tour.is_regular ? "Регулярный тур" : formatDateRange(tour.date_start, tour.date_end)}
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap">
                    {tour.is_regular ? "—" : formatTourDuration(tour.date_start, tour.date_end) || "—"}
                  </td>
                  <td className="px-4 py-4">
                    {tour.slots_left}/{tour.slots_total}
                  </td>
                  <td className="px-4 py-4">
                    {tourShowsPrice(tour) ? formatPrice(tour.price, tour.currency) : "—"}
                  </td>
                  <td className="px-4 py-4">
                    <div className="flex flex-wrap gap-1">
                      <StatusBadge variant={tour.is_active ? "success" : "neutral"}>
                        {tour.is_active ? "Активный" : "Скрыт"}
                      </StatusBadge>
                      {tour.is_regular ? <StatusBadge variant="neutral">Регулярный</StatusBadge> : null}
                      {tour.is_hot ? <StatusBadge variant="warning">Популярный</StatusBadge> : null}
                      {tour.overbooking_enabled ? (
                        <StatusBadge variant="neutral">Овербукинг</StatusBadge>
                      ) : null}
                    </div>
                  </td>
                  <td className="px-4 py-4">
                    <div className="space-y-2">
                      <EditTourForm tour={tour} />
                      <form action={deleteTourAction}>
                        <input type="hidden" name="id" value={tour.id} />
                        <button type="submit" className="btn-danger">
                          Удалить
                        </button>
                      </form>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>

      <CreateTourForm />
    </div>
  );
}
