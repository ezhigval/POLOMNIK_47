"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { bootstrapHomePageAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

export function BootstrapHomeButton() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onClick() {
    setLoading(true);
    setError(null);
    try {
      const page = await bootstrapHomePageAction();
      router.push(`/management/content/${page.id}`);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать главную");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="rounded-2xl border border-brand-200 bg-brand-50/60 p-5">
      <h3 className="font-semibold text-stone-900">Главная страница</h3>
      <p className="mt-1 text-sm text-stone-600">
        Создайте главную с готовыми блоками: шапка, направления, о службе, вопросы и другие секции. Отзывы
        и туры подтягиваются автоматически из своих разделов.
      </p>
      <FormError className="mt-2">{error}</FormError>
      <button
        type="button"
        onClick={onClick}
        disabled={loading}
        className="btn-primary mt-4"
      >
        {loading ? "Создаём…" : "Создать главную из шаблона"}
      </button>
    </div>
  );
}
