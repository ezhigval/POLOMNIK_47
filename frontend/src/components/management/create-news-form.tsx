"use client";

import { FormEvent, useState } from "react";
import { createNewsAction } from "@/app/management/actions";
import { NewsImageField } from "@/components/management/news-image-field";
import { FormError } from "@/components/form-error";

export function CreateNewsForm() {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [formKey, setFormKey] = useState(0);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);

    try {
      await createNewsAction({
        slug: String(formData.get("slug") ?? ""),
        title: String(formData.get("title") ?? ""),
        excerpt: String(formData.get("excerpt") ?? ""),
        body: String(formData.get("body") ?? ""),
        image_url: String(formData.get("image_url") ?? ""),
        published_at: String(formData.get("published_at") ?? ""),
        is_published: formData.get("is_published") === "on",
        is_pinned: formData.get("is_pinned") === "on",
        sort_order: Number(formData.get("sort_order") ?? 0),
      });
      setFormKey((value) => value + 1);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать статью");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form key={formKey} onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <h2 className="text-lg font-semibold">Добавить статью</h2>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Заголовок</span>
        <input required name="title" className="input-field" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Адрес (slug)</span>
        <input required name="slug" className="input-field" placeholder="tikhvin-path" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Дата публикации</span>
        <input required type="date" name="published_at" className="input-field" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Анонс</span>
        <textarea required name="excerpt" rows={2} className="input-field" />
      </label>

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Текст статьи</span>
        <textarea
          required
          name="body"
          rows={8}
          className="input-field"
          placeholder="Абзацы разделяйте пустой строкой. Фотолента: один адрес картинки на строку, без другого текста."
        />
      </label>

      <NewsImageField />

      <label className="block text-sm">
        <span className="mb-1 block font-medium">Порядок среди закреплённых</span>
        <input type="number" name="sort_order" defaultValue={0} className="input-field" />
        <span className="mt-1 block text-xs text-stone-500">
          Меньше число — выше. 1 — главная карточка, 2 и 3 — рядом.
        </span>
      </label>

      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="is_pinned" className="size-4" />
        Закрепить вверху (главная и две рядом)
      </label>

      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="is_published" defaultChecked className="size-4" />
        Опубликовать сразу
      </label>

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading} className="btn-primary">
        {loading ? "Сохраняем..." : "Создать статью"}
      </button>
    </form>
  );
}
