"use client";

import Link from "next/link";

export default function TourError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="mx-auto max-w-3xl px-4 py-20 text-center">
      <h1 className="mb-3 text-2xl font-semibold">Не удалось загрузить тур</h1>
      <p className="mb-6 text-stone-600">
        {error.message || "Попробуйте обновить страницу или вернуться к списку туров."}
      </p>
      <div className="flex flex-wrap items-center justify-center gap-3">
        <button
          type="button"
          onClick={reset}
          className="rounded-full bg-stone-900 px-5 py-2 text-sm font-medium text-white"
        >
          Попробовать снова
        </button>
        <Link
          href="/"
          className="rounded-full border border-stone-300 px-5 py-2 text-sm font-medium text-stone-700"
        >
          К списку туров
        </Link>
      </div>
    </div>
  );
}
