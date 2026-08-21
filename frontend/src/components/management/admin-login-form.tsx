"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

export function AdminLoginForm() {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);
    const response = await fetch("/api/management/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ password: formData.get("password") }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось войти");
      return;
    }

    router.push("/management");
    router.refresh();
  }

  return (
    <form onSubmit={onSubmit} className="mx-auto max-w-sm space-y-4 rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
      <div>
        <h1 className="text-xl font-semibold text-stone-900">Вход в админку</h1>
        <p className="mt-1 text-sm text-stone-600">Введите admin-токен для доступа к управлению</p>
      </div>

      <label className="block text-sm">
        <span className="form-label">Admin token</span>
        <input
          required
          type="password"
          name="password"
          autoComplete="current-password"
          className="input-field"
        />
      </label>

      {error ? <p className="text-sm text-red-700">{error}</p> : null}

      <button type="submit" disabled={loading} className="btn-primary w-full">
        {loading ? "Проверяем..." : "Войти"}
      </button>
    </form>
  );
}
