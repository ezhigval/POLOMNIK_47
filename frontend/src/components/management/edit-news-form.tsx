"use client";

import { FormEvent, useState } from "react";
import { updateNewsAction } from "@/app/management/actions";
import { NewsImageField } from "@/components/management/news-image-field";
import type { ManagementNewsArticle } from "@/lib/api/management";
import { FormError } from "@/components/form-error";

type EditNewsFormProps = {
  article: ManagementNewsArticle;
};

export function EditNewsForm({ article }: EditNewsFormProps) {
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);

    try {
      await updateNewsAction(article.id, {
        slug: String(formData.get("slug") ?? ""),
        title: String(formData.get("title") ?? ""),
        excerpt: String(formData.get("excerpt") ?? ""),
        body: String(formData.get("body") ?? ""),
        image_url: String(formData.get("image_url") ?? ""),
        published_at: String(formData.get("published_at") ?? ""),
        is_published: formData.get("is_published") === "on",
        sort_order: Number(formData.get("sort_order") ?? 0),
      });
      setOpen(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось обновить статью");
    } finally {
      setLoading(false);
    }
  }

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="text-sm font-medium text-brand-800 hover:text-brand-900"
      >
        Редактировать
      </button>
    );
  }

  return (
    <form onSubmit={onSubmit} className="mt-3 space-y-3 rounded-xl border border-stone-200 bg-stone-50 p-3">
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Заголовок</span>
        <input required name="title" defaultValue={article.title} className="input-field" />
      </label>
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Адрес (slug)</span>
        <input required name="slug" defaultValue={article.slug} className="input-field" />
      </label>
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Дата публикации</span>
        <input required type="date" name="published_at" defaultValue={article.published_at} className="input-field" />
      </label>
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Анонс</span>
        <textarea required name="excerpt" rows={2} defaultValue={article.excerpt} className="input-field" />
      </label>
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Текст статьи</span>
        <textarea required name="body" rows={8} defaultValue={article.body} className="input-field" />
      </label>
      <NewsImageField defaultValue={article.image_url} />
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Порядок</span>
        <input type="number" name="sort_order" defaultValue={article.sort_order} className="input-field" />
      </label>
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="is_published" defaultChecked={article.is_published} className="size-4" />
        Опубликована
      </label>
      <FormError>{error}</FormError>
      <div className="flex gap-2">
        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? "Сохраняем..." : "Сохранить"}
        </button>
        <button type="button" onClick={() => setOpen(false)} className="btn-secondary">
          Отмена
        </button>
      </div>
    </form>
  );
}
