"use client";

import { FormEvent, useState } from "react";
import { createSMMPostAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

const CHANNELS = [
  { id: "site_news", label: "site_news" },
  { id: "telegram_channel", label: "telegram_channel" },
  { id: "vk_wall", label: "vk_wall" },
  { id: "max_feed", label: "max_feed" },
];

export function CreateSMMForm() {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [formKey, setFormKey] = useState(0);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    const formData = new FormData(event.currentTarget);
    const channels = CHANNELS.map((item) => item.id).filter((id) => formData.get(`ch_${id}`) === "on");
    const local = String(formData.get("publish_at") ?? "");
    const publishAt = local ? new Date(local).toISOString() : "";
    try {
      await createSMMPostAction({
        title: String(formData.get("title") ?? ""),
        body: String(formData.get("body") ?? ""),
        url: String(formData.get("url") ?? ""),
        publish_at: publishAt,
        channels,
      });
      setFormKey((value) => value + 1);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить материал");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form key={formKey} onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <h2 className="text-lg font-semibold text-stone-900">Новый материал</h2>
      <p className="text-sm text-stone-500">Текст только из этого поля. Каналы publisher не дописывают слоганы.</p>
      <label className="block space-y-1 text-sm">
        <span>Заголовок</span>
        <input name="title" required className="w-full rounded-xl border border-stone-200 px-3 py-2" />
      </label>
      <label className="block space-y-1 text-sm">
        <span>Текст поста</span>
        <textarea name="body" required rows={6} className="w-full rounded-xl border border-stone-200 px-3 py-2" />
      </label>
      <label className="block space-y-1 text-sm">
        <span>URL (необязательно)</span>
        <input name="url" className="w-full rounded-xl border border-stone-200 px-3 py-2" />
      </label>
      <label className="block space-y-1 text-sm">
        <span>Слот времени</span>
        <input name="publish_at" type="datetime-local" required className="w-full rounded-xl border border-stone-200 px-3 py-2" />
      </label>
      <fieldset className="space-y-2 text-sm">
        <legend>Каналы PublisherPort</legend>
        {CHANNELS.map((item) => (
          <label key={item.id} className="flex items-center gap-2">
            <input type="checkbox" name={`ch_${item.id}`} />
            <span>{item.label}</span>
          </label>
        ))}
      </fieldset>
      <FormError>{error}</FormError>
      <button type="submit" disabled={loading} className="btn-primary">
        {loading ? "Сохранение…" : "Сохранить в план"}
      </button>
    </form>
  );
}
