"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { FormError } from "@/components/form-error";
import { HoneypotField } from "@/components/honeypot-field";
import {
  PersonalDataConsentCheckbox,
  PhotoDistributionConsentCheckbox,
} from "@/components/consent-checkbox";

export type UserPhoto = {
  id: string;
  url: string;
  caption: string;
  allow_distribution: boolean;
  created_at: string;
};

type UserPhotoUploadFormProps = {
  photos: UserPhoto[];
};

export function UserPhotoUploadForm({ photos }: UserPhotoUploadFormProps) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [caption, setCaption] = useState("");
  const [consentPersonalData, setConsentPersonalData] = useState(false);
  const [allowDistribution, setAllowDistribution] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);

    if (!consentPersonalData) {
      setLoading(false);
      setError("Необходимо согласие на обработку персональных данных");
      return;
    }

    const form = event.currentTarget;
    const fileInput = form.elements.namedItem("file") as HTMLInputElement | null;
    const file = fileInput?.files?.[0];
    if (!file) {
      setLoading(false);
      setError("Выберите файл изображения");
      return;
    }

    const uploadData = new FormData();
    uploadData.set("file", file);
    const uploadResponse = await fetch("/api/account/uploads", {
      method: "POST",
      body: uploadData,
    });
    if (!uploadResponse.ok) {
      const body = await uploadResponse.json().catch(() => null);
      setLoading(false);
      setError(body?.error ?? "Не удалось загрузить файл");
      return;
    }
    const uploaded = await uploadResponse.json();
    const url = uploaded?.data?.url as string | undefined;
    if (!url) {
      setLoading(false);
      setError("Сервер не вернул URL файла");
      return;
    }

    const response = await fetch("/api/account/photos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        url,
        caption,
        allow_distribution: allowDistribution,
        consent_personal_data: consentPersonalData,
        website: String(new FormData(form).get("website") ?? ""),
      }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось сохранить фотографию");
      return;
    }

    setCaption("");
    setConsentPersonalData(false);
    setAllowDistribution(false);
    form.reset();
    router.refresh();
  }

  async function onDelete(id: string) {
    setLoading(true);
    setError(null);
    const response = await fetch(`/api/account/photos/${id}`, { method: "DELETE" });
    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось удалить фотографию");
      return;
    }
    router.refresh();
  }

  return (
    <div className="space-y-6">
      <form onSubmit={onSubmit} className="relative space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <HoneypotField />
        <div>
          <h2 className="font-display text-xl font-semibold text-stone-900">Загрузить фотографию</h2>
          <p className="mt-1 text-sm text-stone-600">
            JPEG, PNG, WebP или GIF. Публикация на сайте — только с отдельным согласием на распространение.
          </p>
        </div>

        <label className="block text-sm">
          <span className="form-label">Файл</span>
          <input required name="file" type="file" accept="image/*" className="input-field" disabled={loading} />
        </label>

        <label className="block text-sm">
          <span className="form-label">Подпись (необязательно)</span>
          <input
            name="caption"
            className="input-field"
            value={caption}
            onChange={(event) => setCaption(event.target.value)}
            disabled={loading}
          />
        </label>

        <PersonalDataConsentCheckbox
          checked={consentPersonalData}
          onChange={setConsentPersonalData}
          disabled={loading}
        />
        <PhotoDistributionConsentCheckbox
          checked={allowDistribution}
          onChange={setAllowDistribution}
          disabled={loading}
        />

        <FormError>{error}</FormError>

        <button
          type="submit"
          disabled={loading || !consentPersonalData}
          className="btn-primary"
        >
          {loading ? "Загружаем…" : "Загрузить"}
        </button>
      </form>

      <div className="space-y-3">
        <h2 className="font-display text-xl font-semibold text-stone-900">Мои фотографии</h2>
        {photos.length === 0 ? (
          <p className="text-sm text-stone-500">Пока нет загруженных фотографий.</p>
        ) : (
          <ul className="grid gap-4 sm:grid-cols-2">
            {photos.map((photo) => (
              <li key={photo.id} className="overflow-hidden rounded-2xl border border-stone-200 bg-white">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src={photo.url} alt={photo.caption || "Фотография"} className="aspect-[4/3] w-full object-cover" />
                <div className="space-y-2 p-4 text-sm">
                  {photo.caption ? <p className="text-stone-800">{photo.caption}</p> : null}
                  <p className="text-xs text-stone-500">
                    {photo.allow_distribution ? "Разрешена публикация" : "Только в личном кабинете"}
                  </p>
                  <button
                    type="button"
                    className="btn-secondary text-sm"
                    disabled={loading}
                    onClick={() => onDelete(photo.id)}
                  >
                    Удалить
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
