"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useCallback, useState } from "react";
import { PhoneCallVerify } from "@/components/auth/phone-call-verify";
import { FormError } from "@/components/form-error";
import { safeReturnUrl } from "@/lib/site-nav";

type RegisterFormProps = {
  returnUrl?: string;
};

export function RegisterForm({ returnUrl = "/account/trips" }: RegisterFormProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [phone, setPhone] = useState("");
  const [phoneCheckId, setPhoneCheckId] = useState<string | null>(null);
  const [phoneCallAvailable, setPhoneCallAvailable] = useState(false);
  const destination = safeReturnUrl(returnUrl);
  const loginHref =
    destination === "/account/trips"
      ? "/account/login"
      : `/account/login?returnUrl=${encodeURIComponent(destination)}`;

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

    if (phoneCallAvailable && !phoneCheckId) {
      setLoading(false);
      setError("Сначала подтвердите телефон звонком с вашего номера");
      return;
    }

    const formData = new FormData(event.currentTarget);
    const response = await fetch("/api/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: formData.get("name"),
        email: formData.get("email"),
        phone: formData.get("phone"),
        password: formData.get("password"),
        phone_check_id: phoneCheckId,
      }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось зарегистрироваться");
      return;
    }

    router.push(destination);
    router.refresh();
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
      <div>
        <h2 className="font-display text-xl font-semibold text-stone-900">Данные аккаунта</h2>
        <p className="mt-1 text-sm text-stone-600">После регистрации откроется личный кабинет.</p>
      </div>

      <label className="block text-sm">
        <span className="form-label">Имя и фамилия</span>
        <input required name="name" className="input-field" autoComplete="name" />
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
          onChange={(event) => {
            setPhone(event.target.value);
            setPhoneCheckId(null);
          }}
        />
      </label>

      <PhoneCallVerify
        phone={phone}
        onConfirmed={onPhoneConfirmed}
        onAvailable={onPhoneAvailable}
        onUnavailable={onPhoneUnavailable}
      />

      <label className="block text-sm">
        <span className="form-label">Email</span>
        <input name="email" type="email" className="input-field" autoComplete="email" />
      </label>

      <label className="block text-sm">
        <span className="form-label">Пароль</span>
        <input required type="password" name="password" className="input-field" minLength={8} />
      </label>

      <FormError>{error}</FormError>

      <button
        type="submit"
        disabled={loading || (phoneCallAvailable && !phoneCheckId)}
        className="btn-primary w-full"
      >
        {loading ? "Создаём…" : "Создать аккаунт"}
      </button>

      <p className="text-center text-sm text-stone-600">
        Уже есть аккаунт?{" "}
        <Link href={loginHref} className="font-medium text-brand-800 hover:underline">
          Войти
        </Link>
      </p>
    </form>
  );
}
