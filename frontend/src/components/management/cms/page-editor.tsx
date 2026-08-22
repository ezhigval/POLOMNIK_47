"use client";

import { FormEvent, useMemo, useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import {
  createCmsBlockAction,
  deleteCmsBlockAction,
  reorderCmsBlocksAction,
  updateCmsPageAction,
} from "@/app/management/actions";
import { BlockEditor } from "@/components/management/cms/block-editor";
import type { CmsBlockTemplate, CmsPage } from "@/lib/api/cms";
import { FormError } from "@/components/form-error";

type PageEditorProps = {
  page: CmsPage;
  templates: CmsBlockTemplate[];
};

export function PageEditor({ page, templates }: PageEditorProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedType, setSelectedType] = useState(templates[0]?.type ?? "rich_text");
  const [isPending, startTransition] = useTransition();

  const blocks = useMemo(
    () => [...(page.blocks ?? [])].sort((a, b) => a.sort_order - b.sort_order),
    [page.blocks],
  );

  async function onSaveMeta(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    const formData = new FormData(event.currentTarget);

    try {
      await updateCmsPageAction(
        page.id,
        {
          title: String(formData.get("title") ?? ""),
          path: String(formData.get("path") ?? ""),
          is_published: formData.get("is_published") === "on",
        },
        page.slug,
      );
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить страницу");
    } finally {
      setLoading(false);
    }
  }

  async function onAddBlock() {
    setError(null);
    const template = templates.find((item) => item.type === selectedType);
    try {
      await createCmsBlockAction(
        page.id,
        {
          type: selectedType,
          content: template?.content ?? {},
          is_visible: true,
        },
        page.slug,
      );
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось добавить блок");
    }
  }

  function moveBlock(index: number, direction: -1 | 1) {
    const nextIndex = index + direction;
    if (nextIndex < 0 || nextIndex >= blocks.length) return;

    const next = [...blocks];
    const [item] = next.splice(index, 1);
    next.splice(nextIndex, 0, item);
    const blockIds = next.map((block) => block.id);

    startTransition(async () => {
      setError(null);
      try {
        await reorderCmsBlocksAction(page.id, blockIds, page.slug);
        router.refresh();
      } catch (err) {
        setError(err instanceof Error ? err.message : "Не удалось изменить порядок");
      }
    });
  }

  return (
    <div className="space-y-8">
      <form
        onSubmit={onSaveMeta}
        className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5"
      >
        <div>
          <h2 className="text-lg font-semibold">Метаданные</h2>
          <p className="mt-1 text-sm text-stone-600">
            Slug: <span className="font-medium text-stone-900">{page.slug}</span>
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <label className="block text-sm">
            <span className="mb-1 block font-medium">Заголовок</span>
            <input required name="title" defaultValue={page.title} className="input-field" />
          </label>
          <label className="block text-sm">
            <span className="mb-1 block font-medium">Path</span>
            <input required name="path" defaultValue={page.path} className="input-field" />
          </label>
        </div>

        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            name="is_published"
            defaultChecked={page.is_published}
            className="size-4"
          />
          Опубликована
        </label>

        <FormError>{error}</FormError>

        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? "Сохраняем..." : "Сохранить страницу"}
        </button>
      </form>

      <section className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold">Блоки</h2>
            <p className="mt-1 text-sm text-stone-600">{blocks.length} на странице</p>
          </div>

          <div className="flex flex-wrap items-end gap-2">
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Шаблон</span>
              <select
                value={selectedType}
                onChange={(event) => setSelectedType(event.target.value as typeof selectedType)}
                className="input-field min-w-56"
              >
                {templates.map((template) => (
                  <option key={template.type} value={template.type}>
                    {template.label}
                  </option>
                ))}
              </select>
            </label>
            <button type="button" onClick={onAddBlock} className="btn-primary">
              Добавить блок
            </button>
          </div>
        </div>

        {blocks.length === 0 ? (
          <p className="py-8 text-center text-sm text-stone-500">Блоков пока нет.</p>
        ) : (
          <ul className="space-y-4">
            {blocks.map((block, index) => (
              <li
                key={block.id}
                className="rounded-xl border border-stone-200 bg-stone-50/70 p-4"
              >
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="font-medium text-stone-900">
                      #{index + 1} · {block.type}
                    </p>
                    <p className="text-xs text-stone-500">
                      {block.is_visible ? "Видимый" : "Скрыт"} · sort {block.sort_order}
                    </p>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    <button
                      type="button"
                      disabled={isPending || index === 0}
                      onClick={() => moveBlock(index, -1)}
                      className="btn-secondary px-3 py-1.5 text-sm disabled:opacity-40"
                    >
                      ↑
                    </button>
                    <button
                      type="button"
                      disabled={isPending || index === blocks.length - 1}
                      onClick={() => moveBlock(index, 1)}
                      className="btn-secondary px-3 py-1.5 text-sm disabled:opacity-40"
                    >
                      ↓
                    </button>
                    <form action={deleteCmsBlockAction}>
                      <input type="hidden" name="id" value={block.id} />
                      <input type="hidden" name="page_id" value={page.id} />
                      <input type="hidden" name="slug" value={page.slug} />
                      <button type="submit" className="btn-danger">
                        Удалить
                      </button>
                    </form>
                  </div>
                </div>

                <div className="mt-3">
                  <BlockEditor block={block} pageId={page.id} slug={page.slug} />
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
