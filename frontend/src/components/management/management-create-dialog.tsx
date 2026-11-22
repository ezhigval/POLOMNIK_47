"use client";

import { useId, useRef, useState } from "react";
import { CreateNewsForm } from "@/components/management/create-news-form";
import { CreateReviewForm } from "@/components/management/create-review-form";
import { CreateSMMForm } from "@/components/management/create-smm-form";
import { CreateTourForm } from "@/components/management/create-tour-form";
import { PublishLegalDocumentForm } from "@/components/management/publish-legal-document-form";
import type { ManagementTour } from "@/lib/api/management";

type EntryKind = "news" | "smm" | "tours" | "reviews" | "legal";

const labels: Record<EntryKind, string> = {
  news: "Добавить статью",
  smm: "Добавить материал",
  tours: "Добавить тур",
  reviews: "Добавить отзыв",
  legal: "Добавить версию документа",
};

export function ManagementCreateDialog({
  kind,
  tours = [],
}: {
  kind: EntryKind;
  tours?: ManagementTour[];
}) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [opened, setOpened] = useState(false);
  const labelId = useId();

  function open() {
    dialogRef.current?.showModal();
    setOpened(true);
  }

  function close() {
    dialogRef.current?.close();
  }

  return (
    <>
      <button type="button" className="btn-primary shrink-0 px-3 py-2 sm:px-4"
        onClick={open} aria-label={labels[kind]} aria-haspopup="dialog">
        <span aria-hidden="true" className="text-lg leading-none sm:mr-1">+</span>
        <span className="hidden sm:inline">Добавить</span>
      </button>
      <dialog
        ref={dialogRef}
        aria-labelledby={labelId}
        onClose={() => setOpened(false)}
        onClick={(event) => { if (event.target === event.currentTarget) close(); }}
        className="fixed inset-0 m-auto max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-3xl overflow-y-auto rounded-2xl border border-stone-200 bg-white p-0 text-stone-900 shadow-2xl backdrop:bg-stone-950/60"
      >
        <div className="sticky top-0 z-10 flex items-center justify-between gap-4 border-b border-stone-100 bg-white px-5 py-4">
          <h2 id={labelId} className="text-lg font-semibold">{labels[kind]}</h2>
          <button type="button" onClick={close}
            className="inline-flex size-10 shrink-0 items-center justify-center rounded-full border border-stone-200 text-xl hover:bg-stone-50"
            aria-label="Закрыть окно">×</button>
        </div>
        <div className="p-4 sm:p-6">
          {opened ? (
            <>
              {kind === "news" ? <CreateNewsForm onSuccess={close} /> : null}
              {kind === "smm" ? <CreateSMMForm onSuccess={close} /> : null}
              {kind === "tours" ? <CreateTourForm onSuccess={close} /> : null}
              {kind === "reviews" ? <CreateReviewForm tours={tours} onSuccess={close} /> : null}
              {kind === "legal" ? <PublishLegalDocumentForm onSuccess={close} /> : null}
            </>
          ) : null}
        </div>
      </dialog>
    </>
  );
}
