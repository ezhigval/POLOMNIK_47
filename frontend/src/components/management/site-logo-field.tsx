"use client";

import Image from "next/image";
import { useState } from "react";
import { uploadSiteLogoAction } from "@/app/management/actions";
import { defaultLogo, logoSrc } from "@/lib/brand";

async function prepareLogo(file: File) {
  if (file.size > 10 * 1024 * 1024 || !["image/png", "image/jpeg", "image/webp"].includes(file.type)) {
    throw new Error("Выберите PNG, JPEG или WebP размером до 10 МБ.");
  }
  const bitmap = await createImageBitmap(file);
  try {
    const canvas = document.createElement("canvas");
    canvas.width = canvas.height = 384;
    const ctx = canvas.getContext("2d");
    if (!ctx) throw new Error("Браузер не смог подготовить изображение.");
    const scale = Math.min(384 / bitmap.width, 384 / bitmap.height);
    const w = bitmap.width * scale, h = bitmap.height * scale;
    ctx.imageSmoothingQuality = "high";
    ctx.drawImage(bitmap, (384 - w) / 2, (384 - h) / 2, w, h);
    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, "image/png"));
    if (!blob) throw new Error("Не удалось прочитать изображение.");
    const data = new FormData();
    data.set("file", blob, "logo.png");
    return data;
  } finally {
    bitmap.close();
  }
}

export function SiteLogoField({ value, onChange, onBusy }: {
  value: string;
  onChange: (path: string) => void;
  onBusy: (busy: boolean) => void;
}) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  async function select(file: File) {
    setBusy(true);
    onBusy(true);
    setError("");
    try {
      onChange(await uploadSiteLogoAction(await prepareLogo(file)));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Не удалось загрузить логотип.");
    } finally {
      setBusy(false);
      onBusy(false);
    }
  }
  return <fieldset className="flex flex-wrap items-center gap-5 rounded-xl border border-stone-200 p-4" disabled={busy}>
    <legend className="px-1 text-sm font-semibold">Логотип службы</legend>
    <Image unoptimized src={logoSrc(value)} alt="Предпросмотр логотипа службы" width={80} height={80} className="size-20 object-contain" />
    <div className="min-w-0 space-y-2 text-sm">
      <label className="block font-medium">Выбрать логотип
        <input className="mt-1 block max-w-full file:mr-3 file:rounded-lg file:bg-stone-100 file:px-3 file:py-2" type="file"
          accept="image/png,image/jpeg,image/webp" onChange={(event) => {
            const file = event.currentTarget.files?.[0];
            event.currentTarget.value = "";
            if (file) void select(file);
          }} />
      </label>
      <p className="text-stone-600" role="status">{busy ? "Готовим предпросмотр…" : "PNG, JPEG или WebP до 10 МБ. Для публикации нажмите «Сохранить сайт»."}</p>
      {value && value !== defaultLogo && <button type="button" onClick={() => { setError(""); onChange(defaultLogo); }} className="underline">Вернуть исходный логотип</button>}
      {error && <p role="alert" className="text-red-700">{error}</p>}
    </div>
  </fieldset>;
}
