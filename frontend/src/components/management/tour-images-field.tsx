"use client";

import { ChangeEvent, useId, useState } from "react";
import { uploadTourImageAction } from "@/app/management/actions";
import { FormError } from "@/components/form-error";

type TourImagesFieldProps = {
  name?: string;
  defaultValue?: string;
  rows?: number;
};

export function TourImagesField({
  name = "images",
  defaultValue = "",
  rows = 2,
}: TourImagesFieldProps) {
  const inputId = useId();
  const [value, setValue] = useState(defaultValue);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onFileChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) {
      return;
    }

    setUploading(true);
    setError(null);

    const formData = new FormData();
    formData.append("file", file);

    try {
      const uploaded = await uploadTourImageAction(formData);
      setValue((current) => {
        const lines = current
          .split("\n")
          .map((line) => line.trim())
          .filter(Boolean);
        return [...lines, uploaded.url].join("\n");
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось загрузить фото");
    } finally {
      setUploading(false);
    }
  }

  return (
    <div className="space-y-2 md:col-span-2">
      <label className="block text-sm" htmlFor={inputId}>
        <span className="mb-1 block font-medium">Фото (URL, по одному в строке)</span>
        <textarea
          id={inputId}
          name={name}
          rows={rows}
          value={value}
          onChange={(event) => setValue(event.target.value)}
          className="input-field resize-y"
          placeholder="https://example.com/photo.jpg"
        />
      </label>

      <label className="flex flex-wrap items-center gap-3 text-sm">
        <span className="font-medium text-stone-700">Или загрузить файл</span>
        <input
          type="file"
          accept="image/jpeg,image/png,image/webp,image/gif"
          onChange={onFileChange}
          disabled={uploading}
          className="max-w-full text-sm file:mr-3 file:rounded-lg file:border-0 file:bg-stone-900 file:px-3 file:py-1.5 file:text-sm file:font-medium file:text-white hover:file:bg-stone-800 disabled:opacity-60"
        />
        {uploading ? <span className="text-stone-500">Загрузка…</span> : null}
      </label>

      <FormError>{error}</FormError>
    </div>
  );
}
