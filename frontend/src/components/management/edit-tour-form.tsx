"use client";

import { FormEvent, useState } from "react";
import { updateTourAction } from "@/app/management/actions";
import { TourImagesField } from "@/components/management/tour-images-field";
import { TourRegularCheckbox } from "@/components/management/tour-regular-checkbox";
import { parseImageUrls } from "@/lib/parse-image-urls";
import type { ManagementTour } from "@/lib/api/management";
import { FormError } from "@/components/form-error";

type EditTourFormProps = {
  tour: ManagementTour;
};

export function EditTourForm({ tour }: EditTourFormProps) {
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [isRegular, setIsRegular] = useState(Boolean(tour.is_regular));

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);

    try {
      await updateTourAction(tour.id, {
        slug: String(formData.get("slug") ?? ""),
        title: String(formData.get("title") ?? ""),
        description: String(formData.get("description") ?? ""),
        price: isRegular ? 0 : Number(formData.get("price") ?? 0),
        currency: String(formData.get("currency") ?? "RUB") || "RUB",
        date_start: isRegular ? "" : String(formData.get("date_start") ?? ""),
        date_end: isRegular ? "" : String(formData.get("date_end") ?? ""),
        slots_total: Number(formData.get("slots_total") ?? 0),
        slots_left: Number(formData.get("slots_left") ?? 0),
        location: String(formData.get("location") ?? ""),
        images: parseImageUrls(String(formData.get("images") ?? "")),
        is_active: formData.get("is_active") === "on",
        is_hot: formData.get("is_hot") === "on",
        is_regular: isRegular,
        overbooking_enabled: formData.get("overbooking_enabled") === "on",
      });
      setOpen(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось обновить тур");
    } finally {
      setLoading(false);
    }
  }

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="text-sm font-medium text-brand-800 hover:text-brand-900"
      >
        Редактировать
      </button>
    );
  }

  return (
    <form onSubmit={onSubmit} className="mt-3 space-y-3 rounded-xl border border-stone-200 bg-stone-50 p-3">
      <div className="grid gap-3 md:grid-cols-2">
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Slug</span>
          <input required name="slug" defaultValue={tour.slug} className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Название</span>
          <input required name="title" defaultValue={tour.title} className="input-field" />
        </label>
        <label className="block text-sm md:col-span-2">
          <span className="mb-1 block font-medium">Описание</span>
          <textarea required name="description" rows={2} defaultValue={tour.description} className="input-field" />
        </label>
        {isRegular ? (
          <input type="hidden" name="currency" value={tour.currency || "RUB"} />
        ) : (
          <>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Цена</span>
              <input required type="number" min={0} name="price" defaultValue={tour.price ?? 0} className="input-field" />
            </label>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Валюта</span>
              <input required name="currency" defaultValue={tour.currency} className="input-field" />
            </label>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Дата начала</span>
              <input required type="date" name="date_start" defaultValue={tour.date_start ?? ""} className="input-field" />
            </label>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Дата окончания</span>
              <input required type="date" name="date_end" defaultValue={tour.date_end ?? ""} className="input-field" />
            </label>
          </>
        )}
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Всего мест</span>
          <input required type="number" min={0} name="slots_total" defaultValue={tour.slots_total} className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Свободно мест</span>
          <input required type="number" min={0} name="slots_left" defaultValue={tour.slots_left} className="input-field" />
        </label>
        <label className="block text-sm md:col-span-2">
          <span className="mb-1 block font-medium">Локация</span>
          <input required name="location" defaultValue={tour.location} className="input-field" />
        </label>
        <TourImagesField defaultValue={tour.images.join("\n")} />
      </div>

      <div className="flex flex-wrap gap-4 text-sm">
        <label className="flex items-center gap-2">
          <input type="checkbox" name="is_active" defaultChecked={tour.is_active} className="size-4" />
          Активный
        </label>
        <label className="flex items-center gap-2">
          <input type="checkbox" name="is_hot" defaultChecked={tour.is_hot} className="size-4" />
          Популярный
        </label>
        <label className="flex items-center gap-2">
          <input type="checkbox" name="overbooking_enabled" defaultChecked={tour.overbooking_enabled} className="size-4" />
          Овербукинг
        </label>
      </div>

      <TourRegularCheckbox checked={isRegular} onChange={setIsRegular} />

      <FormError>{error}</FormError>

      <div className="flex gap-2">
        <button type="submit" disabled={loading} className="btn-primary px-4 py-1.5 text-sm">
          {loading ? "Сохраняем..." : "Сохранить"}
        </button>
        <button type="button" onClick={() => setOpen(false)} className="btn-secondary px-4 py-1.5 text-sm">
          Отмена
        </button>
      </div>
    </form>
  );
}
