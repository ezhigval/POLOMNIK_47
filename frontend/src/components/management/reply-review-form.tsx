"use client";

import { FormEvent, useState } from "react";
import { replyReviewAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

type ReplyReviewFormProps = {
  reviewId: string;
  initialReply: string;
};

export function ReplyReviewForm({ reviewId, initialReply }: ReplyReviewFormProps) {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function save(company_reply: string) {
    setLoading(true);
    setError(null);
    try {
      await replyReviewAction(reviewId, company_reply);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить ответ");
    } finally {
      setLoading(false);
    }
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    await save(String(formData.get("company_reply") ?? ""));
  }

  return (
    <form onSubmit={onSubmit} className="mt-4 space-y-2 border-t border-stone-100 pt-4">
      <label className="block text-sm">
        <span className="mb-1 block font-medium text-stone-800">Ответ от лица службы</span>
        <textarea
          name="company_reply"
          rows={3}
          defaultValue={initialReply}
          placeholder="Поблагодарите паломника или уточните детали поездки"
          className="input-field"
        />
      </label>

      <FormError>{error}</FormError>

      <div className="flex flex-wrap gap-2">
        <button type="submit" disabled={loading} className="btn-secondary">
          {loading ? "Сохраняем..." : initialReply ? "Обновить ответ" : "Опубликовать ответ"}
        </button>
        {initialReply ? (
          <button
            type="button"
            disabled={loading}
            className="btn-danger"
            onClick={() => {
              void save("");
            }}
          >
            Удалить ответ
          </button>
        ) : null}
      </div>
    </form>
  );
}
