"use client";

import { useId, useRef, useState } from "react";
import type { CmsPage } from "@/lib/api/cms";

const deviceWidths = {
  desktop: { label: "Компьютер", width: "100%" },
  tablet: { label: "Планшет", width: "768px" },
  phone: { label: "Телефон", width: "390px" },
} as const;
type Device = keyof typeof deviceWidths;
type Space = "guide" | "trips";

export function ManagementPagePreview({ page }: { page: CmsPage }) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const titleId = useId();
  const [device, setDevice] = useState<Device>("desktop");
  const [space, setSpace] = useState<Space>("guide");
  const [opened, setOpened] = useState(false);
  const isHome = page.slug === "home";
  const iframeSrc = `/management/content/${encodeURIComponent(page.id)}/preview${isHome ? `?space=${space}` : ""}`;

  return (
    <>
      <button type="button" className="btn-secondary text-sm" aria-haspopup="dialog"
        onClick={() => { dialogRef.current?.showModal(); setOpened(true); }}>
        Предпросмотр страницы
      </button>
      <dialog ref={dialogRef} aria-labelledby={titleId}
        onClose={() => setOpened(false)}
        className="fixed inset-0 m-auto max-h-[calc(100dvh-2rem)] h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-7xl overflow-hidden rounded-2xl border border-stone-200 bg-stone-100 p-0 shadow-2xl backdrop:bg-stone-950/60">
        <div className="flex h-full min-h-0 flex-col">
          <div className="flex shrink-0 flex-wrap items-center justify-between gap-3 border-b border-stone-200 bg-white px-4 py-3">
            <div className="min-w-0">
              <h2 id={titleId} className="truncate text-lg font-semibold text-stone-900">Предпросмотр: {page.title}</h2>
              <p className="text-xs text-stone-500">Реальные компоненты сайта · сохранённая версия, включая черновики</p>
            </div>
            <button type="button" onClick={() => dialogRef.current?.close()}
              aria-label="Закрыть предпросмотр"
              className="flex size-10 items-center justify-center rounded-full border border-stone-200 text-xl text-stone-700 hover:bg-stone-50">
              ×
            </button>
          </div>
          <div className="flex shrink-0 flex-wrap items-center justify-between gap-2 border-b border-stone-200 bg-white px-4 py-2">
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
            ) : <span className="text-xs text-stone-500">Предпросмотр CMS</span>}
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
          <div className="min-h-0 flex-1 overflow-auto p-2 sm:p-4">
            <div className="mx-auto h-full min-h-[400px] overflow-hidden rounded-xl border border-stone-200 bg-white shadow-lg"
              style={{ width: deviceWidths[device].width, minWidth: device === "phone" ? 320 : undefined }}>
              {opened ? (
                <iframe key={iframeSrc} src={iframeSrc} title={`Предпросмотр страницы: ${page.title}`}
                  className="h-full min-h-[400px] w-full border-0"
                  referrerPolicy="same-origin" />
              ) : null}
            </div>
          </div>
        </div>
      </dialog>
    </>
  );
}
