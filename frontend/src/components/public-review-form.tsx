"use client";

import { FormEvent, useState } from "react";
import { FormError } from "@/components/form-error";
import { HoneypotField } from "@/components/honeypot-field";
import {
  DistributionConsentCheckbox,
  PersonalDataConsentCheckbox,
} from "@/components/consent-checkbox";
import type { Tour } from "@/lib/api/tours";

type PublicReviewFormProps = {
  tours: Pick<Tour, "id" | "title">[];
};

export function PublicReviewForm({ tours }: PublicReviewFormProps) {
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [publishedHint, setPublishedHint] = useState(false);
  const [loading, setLoading] = useState(false);
  const [consentPersonalData, setConsentPersonalData] = useState(false);
  const [consentDistribution, setConsentDistribution] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!consentPersonalData) {
      setError("Необходимо согласие на обработку персональных данных");
      return;
    }

    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);
    const allowPublication = consentDistribution;
    const response = await fetch("/api/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        tour_id: String(formData.get("tour_id") ?? ""),
        client_name: String(formData.get("client_name") ?? ""),
        rating: Number(formData.get("rating") ?? 5),
        text: String(formData.get("text") ?? ""),
        website: String(formData.get("website") ?? ""),
        consent_personal_data: consentPersonalData,
        allow_distribution: allowPublication,
      }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось отправить отзыв");
      return;
    }

    setPublishedHint(allowPublication);
    setSuccess(true);
    setConsentPersonalData(false);
    setConsentDistribution(false);
    event.currentTarget.reset();
  }

  if (success) {
    return (
      <div className="rounded-2xl border border-emerald-200 bg-emerald-50 p-5 text-sm text-emerald-950" role="status">
        <p className="font-medium">Спасибо! Отзыв принят.</p>
        <p className="mt-2 text-emerald-900/90">
          {publishedHint
            ? "После модерации он может быть опубликован на сайте."
            : "Отзыв сохранён для внутренней обработки и не будет опубликован как публичный персональный контент без согласия на распространение."}
        </p>
        <button type="button" className="btn-secondary mt-4" onClick={() => setSuccess(false)}>
          Оставить ещё отзыв
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="relative space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <HoneypotField />
      <div>
        <h2 className="font-display text-xl font-semibold text-stone-900">Оставить отзыв</h2>
        <p className="mt-1 text-sm text-stone-600">
          Публикация имени и текста — только при отдельном согласии на распространение.
        </p>
      </div>

      <label className="block text-sm">
        <span className="form-label">Тур</span>
        <select required name="tour_id" defaultValue="" className="input-field">
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
        <span className="form-label">Имя</span>
        <input required name="client_name" className="input-field" autoComplete="name" />
      </label>

      <label className="block text-sm">
        <span className="form-label">Оценка</span>
        <select name="rating" defaultValue={5} className="input-field">
          {[5, 4, 3, 2, 1].map((value) => (
            <option key={value} value={value}>
              {value}
            </option>
          ))}
        </select>
      </label>

      <label className="block text-sm">
        <span className="form-label">Текст отзыва</span>
        <textarea required name="text" rows={4} className="input-field resize-y" />
      </label>

      <PersonalDataConsentCheckbox
        checked={consentPersonalData}
        onChange={setConsentPersonalData}
        disabled={loading}
      />
      <DistributionConsentCheckbox
        checked={consentDistribution}
        onChange={setConsentDistribution}
        disabled={loading}
      />

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading || !consentPersonalData || tours.length === 0} className="btn-primary">
        {loading ? "Отправляем…" : "Отправить отзыв"}
      </button>
    </form>
  );
}
