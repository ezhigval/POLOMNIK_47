"use client";

import { FormEvent, useState } from "react";
import { createReviewAction } from "@/app/management/actions";
import type { ManagementTour } from "@/lib/api/management";
import { FormError } from "@/components/form-error";

type CreateReviewFormProps = {
  tours: ManagementTour[];
};

export function CreateReviewForm({ tours }: CreateReviewFormProps) {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);

    try {
      await createReviewAction({
        tour_id: String(formData.get("tour_id") ?? ""),
        client_name: String(formData.get("client_name") ?? ""),
        rating: Number(formData.get("rating") ?? 5),
        text: String(formData.get("text") ?? ""),
        is_approved: formData.get("is_approved") === "on",
      });
      event.currentTarget.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать отзыв");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <h2 className="text-lg font-semibold">Добавить отзыв</h2>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Тур</span>
        <select
          required
          name="tour_id"
          defaultValue=""
          className="input-field"
        >
          <option value="" disabled>
            Выберите тур
          </option>
          {tours.map((tour) => (
            <option key={tour.id} value={tour.id}>
              {tour.title}
            </option>
          ))}
        </select>
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Имя клиента</span>
        <input required name="client_name" className="input-field" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Рейтинг</span>
        <select name="rating" defaultValue={5} className="input-field">
          {[5, 4, 3, 2, 1].map((value) => (
            <option key={value} value={value}>
              {value}
            </option>
          ))}
        </select>
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Текст</span>
        <textarea required name="text" rows={3} className="input-field" />
      </label>

      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="is_approved" className="size-4" />
        Одобрен сразу
      </label>

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading || tours.length === 0} className="btn-primary">
        {loading ? "Сохраняем..." : "Создать отзыв"}
      </button>
    </form>
  );
}
