"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { FormError } from "@/components/form-error";
import { HoneypotField } from "@/components/honeypot-field";
import type { Passenger, User } from "@/lib/api/auth";

type PassengerFormProps = {
  user: User;
  passenger?: Passenger;
  onCancel?: () => void;
};

export function PassengerForm({ user, passenger, onCancel }: PassengerFormProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [name, setName] = useState(passenger?.name ?? user.name);
  const [phone, setPhone] = useState(passenger?.phone ?? user.phone);
  const [birthDate, setBirthDate] = useState(passenger?.birth_date ?? "");
  const [passport, setPassport] = useState(passenger?.passport ?? "");

  function fillFromProfile() {
    setName(user.name);
    setPhone(user.phone);
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);
    const payload = {
      name,
      phone,
      birth_date: birthDate,
      passport,
      website: formData.get("website"),
    };
    const url = passenger ? `/api/account/passengers/${passenger.id}` : "/api/account/passengers";
    const method = passenger ? "PATCH" : "POST";
    const response = await fetch(url, {
      method,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось сохранить пассажира");
      return;
    }

    router.refresh();
    onCancel?.();
  }

  return (
    <form onSubmit={onSubmit} className="relative space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <HoneypotField />
      <div>
        <h2 className="font-display text-xl font-semibold text-stone-900">
          {passenger ? "Изменить пассажира" : "Новый пассажир"}
        </h2>
        <p className="mt-1 text-sm text-stone-600">ФИО, телефон, дата рождения и паспорт. СНИЛС не нужен.</p>
      </div>

      <button type="button" className="btn-secondary text-sm" onClick={fillFromProfile}>
        Подставить из профиля
      </button>

      <label className="block text-sm">
        <span className="form-label">ФИО</span>
        <input
          required
          name="name"
          className="input-field"
          autoComplete="name"
          value={name}
          onChange={(event) => setName(event.target.value)}
        />
      </label>

      <label className="block text-sm">
        <span className="form-label">Телефон</span>
        <input
          required
          name="phone"
          type="tel"
          className="input-field"
          autoComplete="tel"
          value={phone}
          onChange={(event) => setPhone(event.target.value)}
        />
      </label>

      <label className="block text-sm">
        <span className="form-label">Дата рождения</span>
        <input
          required
          name="birth_date"
          type="date"
          className="input-field"
          value={birthDate}
          onChange={(event) => setBirthDate(event.target.value)}
        />
      </label>

      <label className="block text-sm">
        <span className="form-label">Паспорт</span>
        <input
          required
          name="passport"
          className="input-field"
          autoComplete="off"
          value={passport}
          onChange={(event) => setPassport(event.target.value)}
        />
      </label>

      <FormError>{error}</FormError>

      <div className="flex flex-wrap gap-2">
        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? "Сохраняем…" : "Сохранить"}
        </button>
        {onCancel ? (
          <button type="button" className="btn-secondary" onClick={onCancel}>
            Отмена
          </button>
        ) : null}
      </div>
    </form>
  );
}
