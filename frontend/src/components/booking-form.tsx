"use client";

import Link from "next/link";
import { FormEvent, useMemo, useRef, useState } from "react";
import { ApiError } from "@/lib/api/client";
import { formatBookingStatus, formatPrice, getSlotsAvailability } from "@/lib/format";
import { createBooking, type CreateBookingResult, type Tour } from "@/lib/api/tours";
import { trackBeginCheckout, trackBookingSubmit } from "@/lib/analytics";
import type { BookingProfile } from "@/lib/auth/user-features";
import { HoneypotField } from "@/components/honeypot-field";
import { MarketingConsentCheckbox, PersonalDataConsentCheckbox } from "@/components/consent-checkbox";

type BookingFormProps = {
  tour: Tour;
  profile?: BookingProfile | null;
};

export function BookingForm({ tour, profile = null }: BookingFormProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<CreateBookingResult | null>(null);
  const [peopleCount, setPeopleCount] = useState(1);
  const [consentPersonalData, setConsentPersonalData] = useState(false);
  const [consentMarketing, setConsentMarketing] = useState(false);
  const beginCheckoutSent = useRef(false);

  const soldOut = getSlotsAvailability(tour.slots_left) === "sold_out";
  const estimatedTotal = useMemo(() => tour.price * peopleCount, [tour.price, peopleCount]);

  function onFormFocus() {
    if (beginCheckoutSent.current || soldOut) {
      return;
    }
    beginCheckoutSent.current = true;
    trackBeginCheckout(tour.id);
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (soldOut) return;
    if (!consentPersonalData) {
      setError("Необходимо согласие на обработку персональных данных");
      return;
    }

    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);

    try {
      const result = await createBooking({
        tour_id: tour.id,
        name: String(formData.get("name") ?? ""),
        phone: String(formData.get("phone") ?? ""),
        email: String(formData.get("email") ?? ""),
        people_count: Number(formData.get("people_count") ?? 1),
        comment: String(formData.get("comment") ?? ""),
        website: String(formData.get("website") ?? ""),
        consent_personal_data: consentPersonalData,
        consent_marketing: consentMarketing,
      });
      setSuccess(result);
      trackBookingSubmit(tour.id, peopleCount);
      setPeopleCount(1);
      setConsentPersonalData(false);
      setConsentMarketing(false);
      event.currentTarget.reset();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Не удалось отправить заявку. Попробуйте позже.");
      }
    } finally {
      setLoading(false);
    }
  }

  if (success) {
    return (
      <div
        id="booking-form"
        className="rounded-2xl border border-emerald-200 bg-emerald-50 p-6 text-emerald-950"
        role="status"
      >
        <div className="mb-4 inline-flex size-10 items-center justify-center rounded-full bg-emerald-100 text-lg">
          ✓
        </div>
        <h3 className="mb-2 text-xl font-semibold">Заявка отправлена</h3>
        <p className="mb-4 text-sm leading-6 text-emerald-900/90">
          {formatBookingStatus(success.booking_status)}. Менеджер свяжется с вами в ближайшее время
          для уточнения деталей.
        </p>
        <dl className="space-y-2 rounded-xl bg-white/60 p-4 text-sm">
          <div className="flex justify-between gap-4">
            <dt className="text-emerald-800/80">Номер заявки</dt>
            <dd className="font-mono font-medium">{success.booking_id.slice(0, 8)}…</dd>
          </div>
          <div className="flex justify-between gap-4">
            <dt className="text-emerald-800/80">Ориентировочная сумма</dt>
            <dd className="font-medium">{formatPrice(success.total_price, tour.currency)}</dd>
          </div>
        </dl>
        <div className="mt-4 flex flex-col gap-2 sm:flex-row">
          {profile ? (
            <Link href="/account/trips" className="btn-primary flex-1 text-center">
              Мои поездки
            </Link>
          ) : (
            <Link
              href={`/account/login?returnUrl=${encodeURIComponent("/account/trips")}`}
              className="btn-primary flex-1 text-center"
            >
              Войти и отслеживать заявку
            </Link>
          )}
          <Link href="/search" className="btn-secondary flex-1 text-center">
            Другие туры
          </Link>
        </div>
        <button
          type="button"
          onClick={() => setSuccess(null)}
          className="btn-secondary mt-2 w-full"
        >
          Отправить ещё одну заявку
        </button>
      </div>
    );
  }

  return (
    <form
      id="booking-form"
      onSubmit={onSubmit}
      onFocus={onFormFocus}
      className="relative space-y-4 rounded-2xl border border-stone-200 bg-white p-5 shadow-md ring-1 ring-stone-100 lg:sticky lg:top-24"
    >
      <div>
        <h3 className="font-display text-xl font-semibold text-stone-900">Оставить заявку</h3>
        <p className="mt-1 text-sm text-stone-600">
          Без оплаты на сайте — менеджер перезвонит в течение рабочего дня.
        </p>
        {profile ? (
          <p className="mt-2 rounded-xl bg-brand-50 px-3 py-2 text-xs text-brand-900">
            Данные подставлены из вашего профиля — при необходимости измените их перед отправкой.
          </p>
        ) : null}
      </div>

      <ul className="flex flex-wrap gap-2 text-xs text-stone-600">
        <li className="rounded-full bg-stone-100 px-2.5 py-1">Без предоплаты</li>
        <li className="rounded-full bg-stone-100 px-2.5 py-1">Ответ в тот же день</li>
        <li className="rounded-full bg-stone-100 px-2.5 py-1">Консультация бесплатно</li>
      </ul>

      <div className="rounded-xl bg-stone-50 p-4">
        <p className="text-sm text-stone-600">Ориентировочная стоимость</p>
        <p className="text-2xl font-semibold text-stone-900">
          {formatPrice(estimatedTotal, tour.currency)}
        </p>
        <p className="text-xs text-stone-500">
          {formatPrice(tour.price, tour.currency)} × {peopleCount} чел.
        </p>
      </div>

      {soldOut ? (
        <p role="alert" className="rounded-xl bg-stone-100 px-3 py-2 text-sm text-stone-700">
          К сожалению, мест на этот тур больше нет. Выберите другую дату или тур.
        </p>
      ) : null}

      <HoneypotField />

      <label className="block text-sm">
        <span className="mb-1.5 block font-medium text-stone-700">Имя и фамилия</span>
        <input
          required
          name="name"
          autoComplete="name"
          disabled={soldOut}
          className="input-field"
          placeholder="Иван Иванов"
          defaultValue={profile?.name ?? ""}
        />
      </label>

      <label className="block text-sm">
        <span className="mb-1.5 block font-medium text-stone-700">Телефон</span>
        <input
          required
          name="phone"
          type="tel"
          autoComplete="tel"
          disabled={soldOut}
          className="input-field"
          placeholder="+7 999 000-00-00"
          defaultValue={profile?.phone ?? ""}
        />
      </label>

      <label className="block text-sm">
        <span className="mb-1.5 block font-medium text-stone-700">
          Email <span className="font-normal text-stone-400">(необязательно)</span>
        </span>
        <input
          type="email"
          name="email"
          autoComplete="email"
          disabled={soldOut}
          className="input-field"
          placeholder="mail@example.com"
          defaultValue={profile?.email ?? ""}
        />
      </label>

      <label className="block text-sm">
        <span className="mb-1.5 block font-medium text-stone-700">Количество человек</span>
        <input
          required
          min={1}
          max={Math.max(tour.slots_left, 1)}
          type="number"
          name="people_count"
          value={peopleCount}
          disabled={soldOut}
          onChange={(event) => setPeopleCount(Number(event.target.value) || 1)}
          className="input-field"
        />
      </label>

      <label className="block text-sm">
        <span className="mb-1.5 block font-medium text-stone-700">
          Комментарий <span className="font-normal text-stone-400">(необязательно)</span>
        </span>
        <textarea
          name="comment"
          rows={3}
          disabled={soldOut}
          className="input-field resize-y"
          placeholder="Пожелания, вопросы, особенности группы"
        />
      </label>

      {error ? (
        <p role="alert" aria-live="polite" className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">
          {error}
        </p>
      ) : null}

      <PersonalDataConsentCheckbox
        checked={consentPersonalData}
        onChange={setConsentPersonalData}
        disabled={loading || soldOut}
      />
      <MarketingConsentCheckbox
        checked={consentMarketing}
        onChange={setConsentMarketing}
        disabled={loading || soldOut}
      />

      <button type="submit" disabled={loading || soldOut || !consentPersonalData} className="btn-primary w-full">
        {loading ? "Отправляем…" : soldOut ? "Мест нет" : "Отправить заявку"}
      </button>
    </form>
  );
}
