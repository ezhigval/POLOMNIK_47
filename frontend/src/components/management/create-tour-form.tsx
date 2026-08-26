"use client";

import { FormEvent, useState } from "react";
import { createTourAction } from "@/app/management/actions";
import { TourImagesField } from "@/components/management/tour-images-field";
import { TourRegularCheckbox } from "@/components/management/tour-regular-checkbox";
import { parseImageUrls } from "@/lib/parse-image-urls";
import { FormError } from "@/components/form-error";

export function CreateTourForm() {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [formKey, setFormKey] = useState(0);
  const [isRegular, setIsRegular] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const form = event.currentTarget;
    const formData = new FormData(form);

    try {
      await createTourAction({
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
      form.reset();
      setIsRegular(false);
      setFormKey((value) => value + 1);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать тур");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <h2 className="text-lg font-semibold">Создать тур</h2>

      <div className="grid gap-4 md:grid-cols-2">
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Slug</span>
          <input required name="slug" className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Название</span>
          <input required name="title" className="input-field" />
        </label>
        <label className="block text-sm md:col-span-2">
          <span className="mb-1 block font-medium">Описание</span>
          <textarea required name="description" rows={3} className="input-field" />
        </label>
        {isRegular ? (
          <input type="hidden" name="currency" value="RUB" />
        ) : (
          <>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Цена</span>
              <input required type="number" min={0} name="price" className="input-field" />
            </label>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Валюта</span>
              <input required name="currency" defaultValue="RUB" className="input-field" />
            </label>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Дата начала</span>
              <input required type="date" name="date_start" className="input-field" />
            </label>
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Дата окончания</span>
              <input required type="date" name="date_end" className="input-field" />
            </label>
          </>
        )}
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Всего мест</span>
          <input required type="number" min={0} name="slots_total" className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Свободно мест</span>
          <input required type="number" min={0} name="slots_left" className="input-field" />
        </label>
        <label className="block text-sm md:col-span-2">
          <span className="mb-1 block font-medium">Локация</span>
          <input required name="location" className="input-field" />
        </label>
        <TourImagesField key={formKey} />
      </div>

      <div className="flex flex-wrap gap-4 text-sm">
        <label className="flex items-center gap-2">
          <input type="checkbox" name="is_active" defaultChecked className="size-4" />
          Активный
        </label>
        <label className="flex items-center gap-2">
          <input type="checkbox" name="is_hot" className="size-4" />
          Популярный
        </label>
        <label className="flex items-center gap-2">
          <input type="checkbox" name="overbooking_enabled" className="size-4" />
          Овербукинг
        </label>
      </div>

      <TourRegularCheckbox checked={isRegular} onChange={setIsRegular} />

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading} className="btn-primary">
        {loading ? "Сохраняем..." : "Создать тур"}
      </button>
    </form>
  );
}
