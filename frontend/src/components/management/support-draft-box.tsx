"use client";

import { useState } from "react";
import { requestSupportDraftAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

type SupportDraftBoxProps = {
  threadId: string;
};

export function SupportDraftBox({ threadId }: SupportDraftBoxProps) {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [draft, setDraft] = useState<string>("");
  const [note, setNote] = useState<string>("");
  const [configured, setConfigured] = useState<boolean | null>(null);

  async function onRequest() {
    setLoading(true);
    setError(null);
    try {
      const result = await requestSupportDraftAction(threadId);
      setDraft(result.draft);
      setNote(result.note);
      setConfigured(result.configured);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось получить черновик");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-3 rounded-xl border border-amber-100 bg-amber-50/70 px-4 py-4">
      <p className="text-sm font-medium text-stone-800">Черновик первой линии</p>
      <p className="text-xs leading-5 text-stone-600">
        Не отправляется клиенту. Ответить должен человек. Цены и богословие модель не выдумывает.
      </p>
      <button type="button" disabled={loading} className="btn-secondary" onClick={onRequest}>
        {loading ? "Готовим…" : "Показать черновик"}
      </button>
      <FormError>{error}</FormError>
      {configured === false ? (
        <p className="text-sm text-stone-600">AIPort не настроен (`AI_ADAPTER=noop`) — черновика нет, нужна эскалация человеку.</p>
      ) : null}
      {note ? <p className="text-xs text-stone-500">{note}</p> : null}
      {draft ? (
        <pre className="whitespace-pre-wrap rounded-lg bg-white px-3 py-3 text-sm leading-6 text-stone-800">{draft}</pre>
      ) : null}
    </div>
  );
}
