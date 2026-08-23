"use client";

import { FormEvent, useState } from "react";
import { FormError } from "@/components/form-error";
import { publishLegalDocumentAction } from "@/app/management/actions";

const DOC_TYPES = [
  { value: "privacy_policy", label: "Политика обработки ПД" },
  { value: "personal_data", label: "Согласие на обработку ПД" },
  { value: "distribution", label: "Согласие на распространение" },
  { value: "marketing", label: "Согласие на рекламу" },
  { value: "cookie", label: "Cookie Policy" },
  { value: "terms", label: "Пользовательское соглашение" },
  { value: "offer", label: "Публичная оферта" },
];

export function PublishLegalDocumentForm() {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [ok, setOk] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setOk(false);
    const formData = new FormData(event.currentTarget);
    try {
      await publishLegalDocumentAction({
        type: String(formData.get("type") ?? ""),
        version: String(formData.get("version") ?? ""),
        title: String(formData.get("title") ?? ""),
        content: String(formData.get("content") ?? ""),
      });
      setOk(true);
      event.currentTarget.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось опубликовать");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <h2 className="text-lg font-semibold">Опубликовать новую версию</h2>
      <p className="text-xs text-stone-500">
        Старая активная версия того же типа будет деактивирована. Исторические версии не удаляются.
      </p>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Тип</span>
        <select required name="type" className="input-field" defaultValue="">
          <option value="" disabled>
            Выберите тип
          </option>
          {DOC_TYPES.map((item) => (
            <option key={item.value} value={item.value}>
              {item.label}
            </option>
          ))}
        </select>
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Версия</span>
        <input required name="version" className="input-field" placeholder="1.1" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Заголовок</span>
        <input required name="title" className="input-field" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">HTML-текст</span>
        <textarea required name="content" rows={10} className="input-field font-mono text-xs" />
      </label>

      <FormError>{error}</FormError>
      {ok ? <p className="text-sm text-emerald-700">Версия опубликована.</p> : null}

      <button type="submit" disabled={loading} className="btn-primary">
        {loading ? "Публикуем…" : "Опубликовать"}
      </button>
    </form>
  );
}
