"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useCallback, useState } from "react";
import { PhoneCallVerify } from "@/components/auth/phone-call-verify";
import { FormError } from "@/components/form-error";
import { HoneypotField } from "@/components/honeypot-field";
import type { User } from "@/lib/api/auth";

type ProfileFormProps = {
  user: User;
};

export function ProfileForm({ user }: ProfileFormProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);
  const [loading, setLoading] = useState(false);
  const [name, setName] = useState(user.name);
  const [email, setEmail] = useState(user.email);
  const [phone, setPhone] = useState(user.phone);
  const [phoneCheckId, setPhoneCheckId] = useState<string | null>(null);
  const [phoneCallAvailable, setPhoneCallAvailable] = useState(false);

  const phoneChanged = phone.trim() !== "" && phone.trim() !== user.phone;

  const onPhoneUnavailable = useCallback(() => {
    setPhoneCallAvailable(false);
  }, []);

  const onPhoneAvailable = useCallback(() => {
    setPhoneCallAvailable(true);
  }, []);

  const onPhoneConfirmed = useCallback((checkId: string) => {
    setPhoneCheckId(checkId);
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setSaved(false);

    if (phoneCallAvailable && phoneChanged && !phoneCheckId) {
      setLoading(false);
      setError("Сначала подтвердите телефон звонком с вашего номера");
      return;
    }

    const formData = new FormData(event.currentTarget);
    const response = await fetch("/api/account/profile", {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name,
        email,
        phone,
        phone_check_id: phoneCheckId,
        website: formData.get("website"),
      }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось сохранить профиль");
      return;
    }

    setSaved(true);
    setPhoneCheckId(null);
    router.refresh();
  }

  return (
    <form onSubmit={onSubmit} className="relative space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <HoneypotField />
      <div>
        <h2 className="font-display text-xl font-semibold text-stone-900">Данные аккаунта</h2>
        <p className="mt-1 text-sm text-stone-600">Имя, телефон и почта. UUID нужен для назначения роли в админке.</p>
      </div>

      <label className="block text-sm">
        <span className="form-label">Имя и фамилия</span>
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
          name="phone"
          type="tel"
          className="input-field"
          autoComplete="tel"
          value={phone}
          onChange={(event) => {
            setPhone(event.target.value);
            setPhoneCheckId(null);
          }}
        />
      </label>

      {phoneChanged ? (
        <PhoneCallVerify
          phone={phone}
          onConfirmed={onPhoneConfirmed}
          onAvailable={onPhoneAvailable}
          onUnavailable={onPhoneUnavailable}
        />
      ) : null}

      <label className="block text-sm">
        <span className="form-label">Email</span>
        <input
          name="email"
          type="email"
          className="input-field"
          autoComplete="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
        />
      </label>

      <div>
        <p className="text-xs uppercase tracking-wide text-stone-500">UUID пользователя</p>
        <p className="mt-1 break-all font-mono text-sm text-stone-900">{user.id}</p>
      </div>

      <p className="text-sm text-stone-600">
        Пароль меняется через{" "}
        <Link href="/account/forgot-password" className="font-medium text-brand-800 hover:underline">
          восстановление
        </Link>
        .
      </p>

      <FormError>{error}</FormError>
      {saved ? <p className="text-sm text-emerald-800">Сохранено.</p> : null}

      <button
        type="submit"
        disabled={loading || (phoneCallAvailable && phoneChanged && !phoneCheckId)}
        className="btn-primary"
      >
        {loading ? "Сохраняем…" : "Сохранить"}
      </button>
    </form>
  );
}
