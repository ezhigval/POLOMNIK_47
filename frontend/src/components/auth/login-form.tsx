"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { safeReturnUrl } from "@/lib/site-nav";
import { FormError } from "@/components/form-error";

type LoginFormProps = {
  returnUrl?: string;
};

export function LoginForm({ returnUrl = "/account/trips" }: LoginFormProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const destination = safeReturnUrl(returnUrl);
  const registerHref =
    destination === "/account/trips"
      ? "/account/register"
      : `/account/register?returnUrl=${encodeURIComponent(destination)}`;

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

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
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

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading} className="btn-primary w-full">
        {loading ? "Входим…" : "Войти"}
      </button>

      <p className="text-center text-sm text-stone-600">
        Нет аккаунта?{" "}
        <Link href={registerHref} className="font-medium text-brand-800 hover:underline">
          Зарегистрироваться
        </Link>
      </p>
    </form>
  );
}
