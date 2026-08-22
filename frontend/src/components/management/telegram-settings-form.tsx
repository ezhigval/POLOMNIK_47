"use client";

import { FormEvent, useState } from "react";
import { updateTelegramSettingsAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";
import type { ManagementTelegramSettings } from "@/lib/api/management";

type TelegramSettingsFormProps = {
  settings: ManagementTelegramSettings;
};

export function TelegramSettingsForm({ settings }: TelegramSettingsFormProps) {
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);
  const [loading, setLoading] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setSaved(false);

    const formData = new FormData(event.currentTarget);

    try {
      await updateTelegramSettingsAction({
        booking_usernames: String(formData.get("booking_usernames") ?? ""),
        support_usernames: String(formData.get("support_usernames") ?? ""),
      });
      setSaved(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить получателей Telegram");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-8">
      <section className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <div>
          <h2 className="text-lg font-semibold text-stone-900">Telegram</h2>
          <p className="mt-1 text-sm text-stone-600">
            Один бот на уведомления и поддержку. Укажите username без обязательного @ — по одному в
            строке или через запятую. Человек должен один раз написать боту (/start), иначе сообщение
            ему не уйдёт.
          </p>
        </div>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Новые заявки и смена статуса</span>
          <textarea
            name="booking_usernames"
            rows={3}
            defaultValue={settings.booking_usernames}
            className="input-field"
            placeholder="ezhigval"
          />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Сообщения в поддержку</span>
          <textarea
            name="support_usernames"
            rows={3}
            defaultValue={settings.support_usernames}
            className="input-field"
            placeholder="ezhigval"
          />
        </label>
        {(settings.recipients ?? []).length > 0 ? (
          <ul className="space-y-1 text-sm text-stone-600">
            {(settings.recipients ?? []).map((item) => (
              <li key={`${item.kind}-${item.username}`}>
                @{item.username} · {item.kind === "booking" ? "заявки" : "поддержка"} · {item.status}
              </li>
            ))}
          </ul>
        ) : null}
      </section>

      <div className="flex flex-wrap items-center gap-3">
        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? "Сохраняем…" : "Сохранить"}
        </button>
        {saved ? <p className="text-sm text-emerald-700">Получатели сохранены.</p> : null}
        <FormError>{error}</FormError>
      </div>
    </form>
  );
}
