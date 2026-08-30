"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useCallback, useState } from "react";
import { AuthExpandable } from "@/components/auth/auth-expandable";
import { PhoneCallVerify } from "@/components/auth/phone-call-verify";
import { SocialAuthButtons } from "@/components/auth/social-auth-buttons";
import { FormError } from "@/components/form-error";
import { safeReturnUrl } from "@/lib/site-nav";
import { HoneypotField } from "@/components/honeypot-field";

type LoginFormProps = {
  returnUrl?: string;
};

export function LoginForm({ returnUrl = "/account/trips" }: LoginFormProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [callPhone, setCallPhone] = useState("");
  const [callAvailable, setCallAvailable] = useState(false);
  const destination = safeReturnUrl(returnUrl);
  const registerHref =
    destination === "/account/trips"
      ? "/account/register"
      : `/account/register?returnUrl=${encodeURIComponent(destination)}`;

  const onCallUnavailable = useCallback(() => {
    setCallAvailable(false);
  }, []);

  const onCallAvailable = useCallback(() => {
    setCallAvailable(true);
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        login: formData.get("login"),
        password: formData.get("password"),
        website: formData.get("website"),
      }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось войти");
      return;
    }

    router.push(destination);
    router.refresh();
  }

  async function onCallConfirmed(checkId: string) {
    setLoading(true);
    setError(null);
    const response = await fetch("/api/auth/phone/complete", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ check_id: checkId }),
    });
    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось войти по звонку");
      return;
    }
    router.push(destination);
    router.refresh();
  }

  return (
    <div className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
      <form onSubmit={onSubmit} className="relative space-y-4">
        <HoneypotField />
        <div>
          <h2 className="font-display text-xl font-semibold text-stone-900">Телефон или email</h2>
          <p className="mt-1 text-sm text-stone-600">Введите данные, указанные при регистрации.</p>
        </div>

        <label className="block text-sm">
          <span className="form-label">Телефон или email</span>
          <input required name="login" className="input-field" placeholder="+7 999 000-00-00" />
        </label>

        <label className="block text-sm">
          <span className="form-label">Пароль</span>
          <input required type="password" name="password" className="input-field" minLength={8} />
        </label>

        <p className="text-right text-sm">
          <Link href="/account/forgot-password" className="font-medium text-brand-800 hover:underline">
            Забыли пароль?
          </Link>
        </p>

        <FormError>{error}</FormError>

        <button type="submit" disabled={loading} className="btn-primary w-full">
          {loading ? "Входим…" : "Войти"}
        </button>
      </form>

      <SocialAuthButtons />

      <p className="text-center text-sm text-stone-600">
        Нет аккаунта?{" "}
        <Link href={registerHref} className="font-medium text-brand-800 hover:underline">
          Зарегистрироваться
        </Link>
      </p>

      <AuthExpandable
        title="Вход по звонку"
        hint="Подтверждение звонком с вашего номера, не SMS"
      >
        <label className="block text-sm">
          <span className="form-label">Телефон аккаунта</span>
          <input
            type="tel"
            className="input-field"
            autoComplete="tel"
            value={callPhone}
            onChange={(event) => setCallPhone(event.target.value)}
            placeholder="+7 999 000-00-00"
          />
        </label>
        <PhoneCallVerify
          phone={callPhone}
          disabled={loading}
          onConfirmed={onCallConfirmed}
          onAvailable={onCallAvailable}
          onUnavailable={onCallUnavailable}
        />
        {!callAvailable ? null : (
          <p className="text-xs text-stone-500">Аккаунт с этим телефоном должен уже существовать.</p>
        )}
      </AuthExpandable>
    </div>
  );
}
