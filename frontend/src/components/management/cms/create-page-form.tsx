"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { createCmsPageAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

export function CreatePageForm() {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(event.currentTarget);

    try {
      const page = await createCmsPageAction({
        slug: String(formData.get("slug") ?? ""),
        title: String(formData.get("title") ?? ""),
        path: String(formData.get("path") ?? ""),
        is_published: formData.get("is_published") === "on",
      });
      router.push(`/management/content/${page.id}`);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать страницу");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
      <h2 className="text-lg font-semibold">Создать страницу</h2>

      <div className="grid gap-4">
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Slug</span>
          <input required name="slug" placeholder="home" className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Заголовок</span>
          <input required name="title" placeholder="Главная" className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Path</span>
          <input required name="path" placeholder="/" className="input-field" />
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" name="is_published" defaultChecked className="size-4" />
          Опубликована
        </label>
      </div>

      <FormError>{error}</FormError>

      <button type="submit" disabled={loading} className="btn-primary">
        {loading ? "Создаём..." : "Создать страницу"}
      </button>
    </form>
  );
}
