"use client";

import { useId, useRef, useState } from "react";
import type { CmsBlock, CmsPage } from "@/lib/api/cms";

const deviceWidths = {
  desktop: { label: "Компьютер", width: "100%" },
  tablet: { label: "Планшет", width: "768px" },
  phone: { label: "Телефон", width: "390px" },
} as const;
type Device = keyof typeof deviceWidths;
const guideTypes = new Set(["information_hero", "holy_places", "service_profile", "service_team", "legal_donation", "labor_trips"]);

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
  const [device, setDevice] = useState<Device>("desktop");
  const [space, setSpace] = useState<"guide" | "trips">("guide");
  const isHome = page.slug === "home";
  const blocks = [...(page.blocks ?? [])].filter((block) => block.is_visible).sort((a, b) => a.sort_order - b.sort_order);
  const visibleBlocks = isHome ? blocks.filter((block) => guideTypes.has(block.type) === (space === "guide")) : blocks;
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
            <p className="text-xs text-stone-500">Эскиз сохранённых блоков без публикации. Это не точный рендер сайта; несохранённые изменения не отображаются.</p>
          </div>
          <button type="button" aria-label="Закрыть предпросмотр" onClick={() => ref.current?.close()}
            className="flex size-10 shrink-0 items-center justify-center rounded-full border border-stone-200 text-xl text-stone-700 hover:bg-stone-50">×</button>
        </div>
        <div className="sticky top-[81px] z-10 flex flex-wrap items-center justify-between gap-3 border-b border-stone-200 bg-white px-4 py-3">
          {isHome ? (
            <div role="group" aria-label="Пространство главной" className="flex flex-wrap gap-1">
              {([{ key: "guide", label: "Путеводитель" }, { key: "trips", label: "Подобрать поездку" }] as const).map((item) => (
                <button type="button" key={item.key} aria-pressed={space === item.key}
                  onClick={() => setSpace(item.key)}
                  className={`rounded-lg px-3 py-2 text-xs font-medium ${space === item.key ? "bg-brand-800 text-white" : "bg-stone-100 text-stone-700 hover:bg-stone-200"}`}>
                  {item.label}
                </button>
              ))}
            </div>
          ) : <span className="text-xs text-stone-500">Структура страницы</span>}
          <div role="group" aria-label="Размер экрана предпросмотра" className="flex flex-wrap gap-1">
            {(Object.keys(deviceWidths) as Device[]).map((key) => (
              <button key={key} type="button" aria-pressed={device === key}
                onClick={() => setDevice(key)}
                className={`rounded-lg px-3 py-2 text-xs font-medium ${device === key ? "bg-brand-800 text-white" : "bg-stone-100 text-stone-700 hover:bg-stone-200"}`}>
                {deviceWidths[key].label}
              </button>
            ))}
          </div>
        </div>
        <div className="overflow-x-auto p-3 sm:p-6">
          <div className="mx-auto min-h-80 space-y-4 rounded-xl border border-stone-200 bg-white p-3 shadow-sm transition-[width] sm:p-5"
            style={{ width: deviceWidths[device].width, maxWidth: "100%", minWidth: device === "phone" ? "320px" : undefined }}>
            <p className="text-center text-xs text-stone-400">{deviceWidths[device].label} · {device === "desktop" ? "ширина окна" : deviceWidths[device].width}</p>
            {visibleBlocks.length ? visibleBlocks.map((block) => <PagePreviewBlock key={block.id} block={block} />)
              : <p className="rounded-xl bg-white p-8 text-center text-stone-500">Нет видимых блоков</p>}
          </div>
        </div>
      </dialog>
    </>
  );
}
