"use client";

import { useId, useRef } from "react";
import type { CmsBlock, CmsPage } from "@/lib/api/cms";

function value(content: Record<string, unknown>, key: string) {
  const field = content[key];
  return typeof field === "string" ? field.trim() : "";
}

function PagePreviewBlock({ block }: { block: CmsBlock }) {
  const content = block.content ?? {};
  const title = value(content, "title") || value(content, "eyebrow") || block.type.replaceAll("_", " ");
  const subtitle = value(content, "subtitle") || value(content, "description") || value(content, "lead");
  const body = value(content, "body");
  return (
    <section className="space-y-3 rounded-2xl border border-stone-200 bg-white p-5 sm:p-7">
      <h3 className="font-display text-xl font-semibold text-stone-900">{title}</h3>
      {subtitle ? <p className="text-sm leading-7 text-stone-600">{subtitle}</p> : null}
      {body ? <p className="whitespace-pre-wrap text-sm leading-7 text-stone-700">{body}</p> : null}
      {Array.isArray(content.items) ? (
        <p className="text-xs text-stone-500">Элементов в блоке: {content.items.length}</p>
      ) : null}
    </section>
  );
}

export function ManagementPagePreview({ page }: { page: CmsPage }) {
  const ref = useRef<HTMLDialogElement>(null);
  const id = useId();
  const blocks = [...(page.blocks ?? [])].filter((block) => block.is_visible).sort((a, b) => a.sort_order - b.sort_order);
  return (
    <>
      <button type="button" className="btn-secondary text-sm"
        aria-haspopup="dialog" onClick={() => ref.current?.showModal()}>
        Предпросмотр страницы
      </button>
      <dialog ref={ref} aria-labelledby={id}
        onClick={(event) => { if (event.target === event.currentTarget) ref.current?.close(); }}
        className="fixed inset-0 m-auto max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-5xl overflow-y-auto rounded-2xl bg-stone-50 p-0 shadow-2xl backdrop:bg-stone-950/60">
        <div className="sticky top-0 z-10 flex items-center justify-between gap-3 border-b border-stone-200 bg-white px-5 py-4">
          <div>
            <h2 id={id} className="text-lg font-semibold text-stone-900">Предпросмотр: {page.title}</h2>
            <p className="text-xs text-stone-500">Структура сохранённых блоков, без публикации. Несохранённые правки здесь не видны.</p>
          </div>
          <button type="button" aria-label="Закрыть предпросмотр" onClick={() => ref.current?.close()}
            className="flex size-10 shrink-0 items-center justify-center rounded-full border border-stone-200 text-xl text-stone-700 hover:bg-stone-50">×</button>
        </div>
        <div className="mx-auto max-w-3xl space-y-4 p-4 sm:p-8">
          {blocks.length ? blocks.map((block) => <PagePreviewBlock key={block.id} block={block} />)
            : <p className="rounded-xl bg-white p-8 text-center text-stone-500">Нет видимых блоков</p>}
        </div>
      </dialog>
    </>
  );
}
