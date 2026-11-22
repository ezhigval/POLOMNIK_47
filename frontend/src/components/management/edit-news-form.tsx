"use client";

import { FormEvent, useId, useRef, useState } from "react";
import { updateNewsAction } from "@/app/management/actions";
import { NewsImageField } from "@/components/management/news-image-field";
import type { ManagementNewsArticle } from "@/lib/api/management";
import { FormError } from "@/components/form-error";
import { ManagementFormPreview } from "@/components/management/management-form-preview";

type EditNewsFormProps = {
  article: ManagementNewsArticle;
};

export function EditNewsForm({ article }: EditNewsFormProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const titleId = useId();
  const [open, setOpen] = useState(false);
  function close() {
    dialogRef.current?.close();
  }
  function openDialog() {
    dialogRef.current?.showModal();
    setOpen(true);
  }
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
        is_pinned: formData.get("is_pinned") === "on",
        sort_order: Number(formData.get("sort_order") ?? 0),
      });
      close();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось обновить статью");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <button type="button" onClick={openDialog}
        aria-haspopup="dialog"
        className="text-sm font-medium text-brand-800 hover:text-brand-900">
        Редактировать
      </button>
      <dialog ref={dialogRef} aria-labelledby={titleId}
        onClose={() => { setOpen(false); setError(null); }}
        onClick={(event) => { if (event.target === event.currentTarget && !loading) close(); }}
        className="fixed inset-0 m-auto max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-3xl overflow-y-auto rounded-2xl border border-stone-200 bg-white p-0 shadow-2xl backdrop:bg-stone-950/60">
        <header className="sticky top-0 z-10 flex items-center justify-between gap-3 border-b border-stone-100 bg-white px-5 py-4">
          <div className="min-w-0">
            <h2 id={titleId} className="text-lg font-semibold text-stone-900">
              Редактировать новость
            </h2>
            <p className="truncate text-sm text-stone-500">{article.title}</p>
          </div>
          <button type="button" onClick={close} disabled={loading}
            aria-label="Закрыть редактор"
            className="flex size-10 shrink-0 items-center justify-center rounded-full border border-stone-200 text-xl text-stone-700 hover:bg-stone-50 disabled:opacity-50">
            ×
          </button>
        </header>
        {open ? (
    <form onSubmit={onSubmit} className="space-y-4 p-4 sm:p-6">
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
        <textarea
          required
          name="body"
          rows={8}
          defaultValue={article.body}
          className="input-field"
          placeholder="Абзацы разделяйте пустой строкой. Фотолента: один адрес картинки на строку, без другого текста."
        />
      </label>
      <NewsImageField defaultValue={article.image_url} />
      <label className="block text-sm">
        <span className="mb-1 block font-medium">Порядок среди закреплённых</span>
        <input type="number" name="sort_order" defaultValue={article.sort_order} className="input-field" />
        <span className="mt-1 block text-xs text-stone-500">
          Меньше число — выше. 1 — главная карточка, 2 и 3 — рядом.
        </span>
      </label>
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="is_pinned" defaultChecked={article.is_pinned} className="size-4" />
        Закрепить вверху (главная и две рядом)
      </label>
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="is_published" defaultChecked={article.is_published} className="size-4" />
        Опубликована
      </label>
      <FormError>{error}</FormError>
      <div className="flex gap-2">
        <ManagementFormPreview kind="news" />
      <button type="submit" disabled={loading} className="btn-primary">
          {loading ? "Сохраняем..." : "Сохранить"}
        </button>
        <button type="button" onClick={close} disabled={loading} className="btn-secondary">
          Отмена
        </button>
      </div>
    </form>
        ) : null}
      </dialog>
    </>
  );
}
