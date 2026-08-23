"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";
import { FormError } from "@/components/form-error";

type ResetPasswordFormProps = {
  token: string;
};

export function ResetPasswordForm({ token }: ResetPasswordFormProps) {
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState(false);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    const formData = new FormData(event.currentTarget);
    const password = String(formData.get("password") ?? "");
    const confirm = String(formData.get("confirm") ?? "");
    if (password !== confirm) {
      setLoading(false);
      setError("Пароли не совпадают");
      return;
    }

    const response = await fetch("/api/auth/reset-password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token, new_password: password }),
    });
    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось обновить пароль");
      return;
    }
    setDone(true);
  }

  if (!token) {
    return (
      <div className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
        <h2 className="font-display text-xl font-semibold text-stone-900">Ссылка неполная</h2>
        <p className="text-sm text-stone-600">Запросите новое письмо для восстановления пароля.</p>
        <Link href="/account/forgot-password" className="btn-secondary inline-flex">
          Запросить ссылку
        </Link>
      </div>
    );
  }

  if (done) {
    return (
      <div className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
        <h2 className="font-display text-xl font-semibold text-stone-900">Пароль обновлён</h2>
        <p className="text-sm text-stone-600">Войдите с новым паролем.</p>
        <Link href="/account/login" className="btn-primary inline-flex">
          Войти
        </Link>
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
      <div>
        <h2 className="font-display text-xl font-semibold text-stone-900">Новый пароль</h2>
        <p className="mt-1 text-sm text-stone-600">Не меньше 8 символов.</p>
      </div>

      <label className="block text-sm">
        <span className="form-label">Пароль</span>
        <input required type="password" name="password" className="input-field" minLength={8} autoComplete="new-password" />
      </label>
      <label className="block text-sm">
        <span className="form-label">Повтор пароля</span>
        <input required type="password" name="confirm" className="input-field" minLength={8} autoComplete="new-password" />
      </label>

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading} className="btn-primary w-full">
        {loading ? "Сохраняем…" : "Сохранить пароль"}
      </button>
    </form>
  );
}
