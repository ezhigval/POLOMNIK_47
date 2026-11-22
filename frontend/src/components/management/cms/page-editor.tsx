"use client";

import { FormEvent, useCallback, useEffect, useMemo, useRef, useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import {
  createCmsBlockAction,
  deleteCmsBlockAction,
  reorderCmsBlocksAction,
  updateCmsPageAction,
} from "@/app/management/actions";
import { BlockEditor } from "@/components/management/cms/block-editor";
import type { CmsBlockTemplate, CmsBlockType, CmsPage } from "@/lib/api/cms";
import { FormError } from "@/components/form-error";

const guideBlockTypes = new Set<CmsBlockType>([
  "information_hero", "holy_places", "service_profile", "service_team",
  "legal_donation", "labor_trips",
]);

type HomeSpace = "guide" | "trips";

type PageEditorProps = {
  page: CmsPage;
  templates: CmsBlockTemplate[];
};

export function PageEditor({ page, templates }: PageEditorProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [metaDirty, setMetaDirty] = useState(false);
  const dirtyBlockIds = useRef(new Set<string>());
  const [selectedType, setSelectedType] = useState(templates[0]?.type ?? "rich_text");
  const [homeSpace, setHomeSpace] = useState<HomeSpace>("guide");
  const isHome = page.slug === "home";
  const [isPending, startTransition] = useTransition();

  const trackDirtyBlock = useCallback((blockId: string, dirty: boolean) => {
    if (dirty) dirtyBlockIds.current.add(blockId);
    else dirtyBlockIds.current.delete(blockId);
  }, []);

  useEffect(() => {
    if (!metaDirty) return;
    function warn(event: BeforeUnloadEvent) {
      event.preventDefault();
      event.returnValue = "";
    }
    window.addEventListener("beforeunload", warn);
    return () => window.removeEventListener("beforeunload", warn);
  }, [metaDirty]);

  function selectHomeSpace(value: HomeSpace) {
    if (value === homeSpace) return;
    if ((metaDirty || dirtyBlockIds.current.size > 0) &&
      !window.confirm("Есть несохранённые изменения. Переключить пространство без сохранения?")) return;
    setHomeSpace(value);
    setError(null);
  }

  const blocks = useMemo(
    () => [...(page.blocks ?? [])].sort((a, b) => a.sort_order - b.sort_order),
    [page.blocks],
  );

  const visibleBlocks = useMemo(
    () => isHome
      ? blocks.filter((block) => guideBlockTypes.has(block.type) === (homeSpace === "guide"))
      : blocks,
    [blocks, isHome, homeSpace],
  );
  const visibleTemplates = useMemo(
    () => isHome
      ? templates.filter((template) => guideBlockTypes.has(template.type) === (homeSpace === "guide"))
      : templates,
    [templates, isHome, homeSpace],
  );
  const activeType = visibleTemplates.some((template) => template.type === selectedType)
    ? selectedType
    : (visibleTemplates[0]?.type ?? "");

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
      setMetaDirty(false);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить страницу");
    } finally {
      setLoading(false);
    }
  }

  async function onAddBlock() {
    setError(null);
    const template = visibleTemplates.find((item) => item.type === activeType);
    if (!template) return;
    try {
      await createCmsBlockAction(
        page.id,
        {
          type: template.type,
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
    if (nextIndex < 0 || nextIndex >= visibleBlocks.length) return;

    const orderedIds = visibleBlocks.map((block) => block.id);
    const [item] = orderedIds.splice(index, 1);
    orderedIds.splice(nextIndex, 0, item);

    // Preserve the relative positions of blocks belonging to the other home space.
    const visibleIds = new Set(orderedIds);
    let cursor = 0;
    const blockIds = blocks.map((block) =>
      visibleIds.has(block.id) ? orderedIds[cursor++] : block.id,
    );

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
        onChangeCapture={() => setMetaDirty(true)}
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
        {isHome ? (
          <div className="space-y-3 border-b border-stone-100 pb-5">
            <div>
              <h2 className="text-lg font-semibold text-stone-900">Пространство главной</h2>
              <p className="mt-1 text-sm text-stone-600">
                Выберите пространство для настройки блоков. На сайте переключатель работает так же.
              </p>
            </div>
            <div role="group" aria-label="Пространство главной страницы"
              className="grid w-full grid-cols-2 gap-1 rounded-2xl bg-stone-100 p-1 sm:max-w-xl">
              {([
                { value: "guide", label: "Путеводитель" },
                { value: "trips", label: "Подобрать поездку" },
              ] as const).map(({ value, label }) => (
                <button key={value} type="button"
                  aria-pressed={homeSpace === value}
                  onClick={() => selectHomeSpace(value)}
                  className={`min-h-11 rounded-xl px-3 py-2 text-sm font-semibold transition focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-700 ${homeSpace === value
                    ? "bg-brand-800 text-white shadow-sm"
                    : "text-stone-700 hover:bg-white/80"}`}>
                  {label}
                </button>
              ))}
            </div>
          </div>
        ) : null}
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold">Блоки</h2>
            <p className="mt-1 text-sm text-stone-600">
              {isHome ? `${visibleBlocks.length} в выбранном пространстве` : `${blocks.length} на странице`}
            </p>
          </div>

          <div className="flex flex-wrap items-end gap-2">
            <label className="block text-sm">
              <span className="mb-1 block font-medium">Шаблон</span>
              <select
                value={activeType}
                onChange={(event) => setSelectedType(event.target.value as CmsBlockType)}
                className="input-field min-w-56"
              >
                {visibleTemplates.map((template) => (
                  <option key={template.type} value={template.type}>
                    {template.label}
                  </option>
                ))}
              </select>
            </label>
            <button type="button" onClick={onAddBlock} disabled={!activeType} className="btn-primary disabled:opacity-50">
              Добавить блок
            </button>
          </div>
        </div>

        {visibleBlocks.length === 0 ? (
          <p className="py-8 text-center text-sm text-stone-500">
            {isHome ? "В этом пространстве пока нет блоков." : "Блоков пока нет."}
          </p>
        ) : (
          <ul className="space-y-4">
            {visibleBlocks.map((block, index) => (
              <li
                key={block.id}
                className="rounded-xl border border-stone-200 bg-stone-50/70 p-4"
              >
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="font-medium text-stone-900">
                      {templates.find((item) => item.type === block.type)?.label ?? block.type}
                    </p>
                    <p className="text-xs text-stone-500">
                      {block.is_visible ? "На сайте" : "Скрыт"} · Блок {index + 1}
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
                      disabled={isPending || index === visibleBlocks.length - 1}
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
                  <BlockEditor block={block} pageId={page.id} slug={page.slug} onDirtyChange={trackDirtyBlock} />
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
