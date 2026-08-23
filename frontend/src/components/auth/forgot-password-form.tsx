"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { FormError } from "@/components/form-error";

type MailMethod = {
  available: boolean;
  message?: string;
};

export function ForgotPasswordForm() {
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState(false);
  const [loading, setLoading] = useState(false);
  const [mail, setMail] = useState<MailMethod | null>(null);

  useEffect(() => {
    let cancelled = false;
    void (async () => {
      try {
        const response = await fetch("/api/v1/auth/methods", { cache: "no-store" });
        const body = await response.json().catch(() => null);
        if (!cancelled) {
          setMail(
            body?.data?.mail ?? {
              available: false,
              message: "Пока что недоступно, используйте другой вариант.",
            },
          );
        }
      } catch {
        if (!cancelled) {
          setMail({ available: false, message: "Пока что недоступно, используйте другой вариант." });
        }
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (mail && !mail.available) {
      setError(mail.message || "Пока что недоступно, используйте другой вариант.");
      return;
    }
    setLoading(true);
    setError(null);
    const formData = new FormData(event.currentTarget);
    const response = await fetch("/api/auth/forgot-password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: formData.get("email") }),
    });
    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось отправить письмо");
      return;
    }
    setDone(true);
  }

  if (done) {
    return (
      <div className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
        <h2 className="font-display text-xl font-semibold text-stone-900">Проверьте почту</h2>
        <p className="text-sm leading-6 text-stone-600">
          Если аккаунт с таким email есть, мы отправили ссылку для восстановления пароля.
        </p>
        <Link href="/account/login" className="btn-secondary inline-flex">
          К входу
        </Link>
      </div>
    );
  }

  const unavailable = mail !== null && !mail.available;

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
      <div>
        <h2 className="font-display text-xl font-semibold text-stone-900">Восстановление пароля</h2>
        <p className="mt-1 text-sm text-stone-600">Укажите email, указанный при регистрации.</p>
      </div>

      {unavailable ? (
        <p className="rounded-xl bg-amber-50 px-3 py-2 text-sm text-amber-900">
          {mail?.message || "Пока что недоступно, используйте другой вариант."}
        </p>
      ) : null}

      <label className="block text-sm">
        <span className="form-label">Email</span>
        <input
          required
          type="email"
          name="email"
          className="input-field"
          autoComplete="email"
          disabled={unavailable}
        />
      </label>

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading || unavailable} className="btn-primary w-full">
        {loading ? "Отправляем…" : "Отправить ссылку"}
      </button>

      <p className="text-center text-sm text-stone-600">
        <Link href="/account/login" className="font-medium text-brand-800 hover:underline">
          Вернуться ко входу
        </Link>
      </p>
    </form>
  );
}
