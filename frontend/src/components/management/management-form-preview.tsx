"use client";

import { useId, useRef, useState } from "react";

type PreviewKind = "news" | "tour" | "review" | "smm" | "legal";
type Snapshot = Record<string, string>;

const titles: Record<PreviewKind, string> = {
  news: "Предпросмотр новости",
  tour: "Предпросмотр тура",
  review: "Предпросмотр отзыва",
  smm: "Предпросмотр публикации",
  legal: "Предпросмотр документа",
};

function text(data: Snapshot, key: string, fallback = "Не заполнено") {
  return data[key]?.trim() || fallback;
}

function PreviewBody({ kind, data }: { kind: PreviewKind; data: Snapshot }) {
  if (kind === "review") {
    const rating = Math.max(1, Math.min(5, Number(data.rating) || 5));
    return (
      <article className="space-y-3 rounded-2xl border border-stone-200 p-5">
        <p className="font-semibold text-stone-900">{text(data, "client_name")}</p>
        <p className="text-amber-600" aria-label={`Оценка: ${rating} из 5`}>
          {"★".repeat(rating)}<span className="text-stone-300">{"★".repeat(5 - rating)}</span>
        </p>
        <p className="whitespace-pre-wrap text-sm leading-7 text-stone-700">{text(data, "text")}</p>
        <p className="text-xs text-stone-500">
          {data.is_approved === "on" && data.allow_distribution === "on"
            ? "Предполагаемый статус: одобрен к публикации"
            : "Предполагаемый статус: не опубликован"}
        </p>
      </article>
    );
  }

  if (kind === "tour") {
    return (
      <article className="overflow-hidden rounded-2xl border border-stone-200 bg-white">
        <div className="space-y-3 p-5 sm:p-7">
          <p className="text-xs font-semibold uppercase tracking-wider text-brand-800">Паломническая поездка</p>
          <h3 className="font-display text-2xl font-semibold text-stone-900">{text(data, "title")}</h3>
          <p className="text-sm text-stone-600">{text(data, "location")} · {data.is_regular === "on" ? "Регулярная" : [data.date_start, data.date_end].filter(Boolean).join(" — ") || "Даты не указаны"}</p>
          <p className="whitespace-pre-wrap text-sm leading-7 text-stone-700">{text(data, "description")}</p>
          <div className="flex flex-wrap gap-3 border-t border-stone-100 pt-4 text-sm">
            <strong>{data.is_regular === "on" ? "Стоимость уточняется" : `${data.price || "0"} ${data.currency || "RUB"}`}</strong>
            <span className="text-stone-600">Места: {data.slots_left || "0"} / {data.slots_total || "0"}</span>
          </div>
        </div>
      </article>
    );
  }

  if (kind === "legal") {
    return (
      <article className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <h3 className="text-xl font-semibold">{text(data, "title")}</h3>
        <p className="text-xs text-stone-500">Версия {text(data, "version")} · Тип: {text(data, "type")}</p>
        <div className="max-h-80 overflow-y-auto whitespace-pre-wrap break-words rounded-xl bg-stone-50 p-4 font-mono text-xs text-stone-700">
          {text(data, "content")}
        </div>
        <p className="text-xs text-stone-500">HTML отображается как исходный текст, без выполнения кода.</p>
      </article>
    );
  }

  return (
    <article className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5 sm:p-7">
      <p className="text-xs font-medium uppercase tracking-widest text-brand-800">
        {kind === "news" ? "События службы" : "Контент-план"}
      </p>
      <h3 className="font-display text-2xl font-semibold leading-tight text-stone-900">{text(data, "title")}</h3>
      {kind === "news" ? (
        <p className="text-xs text-stone-500">{text(data, "published_at", "Дата публикации не указана")}</p>
      ) : null}
      {kind === "news" ? (
        <p className="text-base font-medium leading-7 text-stone-700">{text(data, "excerpt")}</p>
      ) : null}
      <div className="whitespace-pre-wrap break-words text-sm leading-7 text-stone-700">
        {text(data, "body")}
      </div>
      {kind === "smm" && data.url ? <p className="break-all text-sm text-brand-800">{data.url}</p> : null}
    </article>
  );
}

export function ManagementFormPreview({ kind }: { kind: PreviewKind }) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const titleId = useId();
  const [snapshot, setSnapshot] = useState<Snapshot>({});

  function show(event: React.MouseEvent<HTMLButtonElement>) {
    const form = event.currentTarget.closest("form");
    if (!form) return;
    const formData = new FormData(form);
    const next: Snapshot = {};
    for (const [key, value] of formData.entries()) {
      if (typeof value === "string") next[key] = value;
    }
    // Read the regular-tour checkbox even if the surrounding form manages it separately.
    if (kind === "tour" && form.querySelector<HTMLInputElement>('input[name="is_regular"]')?.checked) {
      next.is_regular = "on";
    }
    setSnapshot(next);
    dialogRef.current?.showModal();
  }

  return (
    <>
      <button type="button" onClick={show}
        className="btn-secondary inline-flex min-h-10 items-center gap-2 px-4"
        aria-haspopup="dialog">
        <span aria-hidden="true">◉</span> Предпросмотр
      </button>
      <dialog ref={dialogRef} aria-labelledby={titleId}
        className="fixed inset-0 m-auto max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-3xl overflow-y-auto rounded-2xl border border-stone-200 bg-white p-0 shadow-2xl backdrop:bg-stone-950/60"
        onClick={(event) => { if (event.target === event.currentTarget) dialogRef.current?.close(); }}>
        <div className="sticky top-0 z-10 flex items-center justify-between gap-3 border-b border-stone-100 bg-white px-5 py-4">
          <div>
            <h2 id={titleId} className="text-lg font-semibold text-stone-900">{titles[kind]}</h2>
            <p className="text-xs text-stone-500">Предварительный эскиз без сохранения и публикации</p>
          </div>
          <button type="button" onClick={() => dialogRef.current?.close()}
            aria-label="Закрыть предпросмотр"
            className="flex size-10 shrink-0 items-center justify-center rounded-full border border-stone-200 text-xl text-stone-700 hover:bg-stone-50">
            ×
          </button>
        </div>
        <div className="bg-stone-50 p-4 sm:p-6"><PreviewBody kind={kind} data={snapshot}/></div>
      </dialog>
    </>
  );
}
