"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { FormError } from "@/components/form-error";
import { MarketingConsentCheckbox } from "@/components/consent-checkbox";

type ConsentItem = {
  id: string;
  consent_type: string;
  document_version: string;
  accepted_at: string;
};

const typeLabels: Record<string, string> = {
  personal_data: "Обработка персональных данных",
  marketing: "Рекламные сообщения",
  marketing_revoked: "Отказ от рекламных сообщений",
  distribution: "Распространение ПД",
  distribution_revoked: "Прекращение распространения ПД",
  cookie_all: "Cookie: принять все",
  cookie_essential: "Cookie: только необходимые",
  cookie_reject: "Cookie: отказ",
};

export function AccountConsentsPanel() {
  const [items, setItems] = useState<ConsentItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [marketing, setMarketing] = useState(false);
  const [saving, setSaving] = useState(false);

  async function load() {
    setLoading(true);
    const response = await fetch("/api/account/consents", { cache: "no-store" });
    setLoading(false);
    if (!response.ok) {
      setError("Не удалось загрузить историю согласий");
      return;
    }
    const body = await response.json();
    const list = (body?.data ?? []) as ConsentItem[];
    setItems(list);
    const latestMarketing = list.find(
      (item) => item.consent_type === "marketing" || item.consent_type === "marketing_revoked",
    );
    setMarketing(latestMarketing?.consent_type === "marketing");
  }

  useEffect(() => {
    void Promise.resolve().then(() => {
      void load();
    });
  }, []);

  async function saveMarketing() {
    setSaving(true);
    setError(null);
    const response = await fetch("/api/consents", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        consent_type: marketing ? "marketing" : "marketing_revoked",
      }),
    });
    setSaving(false);
    if (!response.ok) {
      setError("Не удалось сохранить предпочтение");
      return;
    }
    await load();
  }

  async function revokeDistribution() {
    setSaving(true);
    setError(null);
    const response = await fetch("/api/consents", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ consent_type: "distribution_revoked" }),
    });
    setSaving(false);
    if (!response.ok) {
      setError("Не удалось отозвать согласие на распространение");
      return;
    }
    await load();
  }

  return (
    <div className="space-y-6">
      <section className="rounded-2xl border border-stone-200 bg-white p-5">
        <h2 className="font-display text-lg font-semibold text-stone-900">Рекламные сообщения</h2>
        <p className="mt-1 text-sm text-stone-600">
          Необязательно. Отзыв рекламного согласия не удаляет аккаунт.
        </p>
        <div className="mt-4 space-y-4">
          <MarketingConsentCheckbox checked={marketing} onChange={setMarketing} disabled={saving} />
          <button type="button" className="btn-primary" disabled={saving} onClick={() => void saveMarketing()}>
            {saving ? "Сохраняем…" : "Сохранить"}
          </button>
        </div>
      </section>

      <section className="rounded-2xl border border-stone-200 bg-white p-5">
        <h2 className="font-display text-lg font-semibold text-stone-900">Распространение ПД</h2>
        <p className="mt-1 text-sm text-stone-600">
          Можно потребовать прекращения публикации отзывов и фотографий. См.{" "}
          <Link href="/legal/distribution-consent" className="text-brand-800 underline underline-offset-2">
            согласие на распространение
          </Link>
          .
        </p>
        <button type="button" className="btn-secondary mt-4" disabled={saving} onClick={() => void revokeDistribution()}>
          Прекратить распространение
        </button>
      </section>

      <section className="rounded-2xl border border-stone-200 bg-white p-5">
        <h2 className="font-display text-lg font-semibold text-stone-900">История согласий</h2>
        <FormError>{error}</FormError>
        {loading ? (
          <p className="mt-3 text-sm text-stone-500">Загрузка…</p>
        ) : items.length === 0 ? (
          <p className="mt-3 text-sm text-stone-500">Пока нет зафиксированных согласий.</p>
        ) : (
          <ul className="mt-4 divide-y divide-stone-100">
            {items.map((item) => (
              <li key={item.id} className="flex flex-col gap-1 py-3 text-sm sm:flex-row sm:justify-between">
                <span className="font-medium text-stone-800">
                  {typeLabels[item.consent_type] ?? item.consent_type}
                </span>
                <span className="text-stone-500">
                  v{item.document_version} · {new Date(item.accepted_at).toLocaleString("ru-RU")}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
