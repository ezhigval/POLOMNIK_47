"use client";

import { FormEvent, useState } from "react";
import { replySupportAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

type ReplySupportFormProps = {
  threadId: string;
};

export function ReplySupportForm({ threadId }: ReplySupportFormProps) {
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const formData = new FormData(form);
    const body = String(formData.get("body") ?? "").trim();
    if (!body) {
      setError("Введите текст ответа");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      await replySupportAction(threadId, body);
      form.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось отправить ответ");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      <label className="block text-sm">
        <span className="mb-1 block font-medium text-stone-800">Ответ паломнику</span>
        <textarea
          name="body"
          rows={4}
          required
          maxLength={4000}
          placeholder="Напишите ответ — он появится в чате поддержки"
          className="input-field"
        />
      </label>
      <FormError>{error}</FormError>
      <button type="submit" disabled={loading} className="btn-primary">
        {loading ? "Отправляем…" : "Отправить ответ"}
      </button>
    </form>
  );
}
