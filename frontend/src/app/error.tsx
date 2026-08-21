"use client";

import Link from "next/link";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="mx-auto max-w-lg px-4 py-20 text-center">
      <div className="mx-auto mb-4 flex size-14 items-center justify-center rounded-full bg-red-50 text-2xl text-red-600">
        !
      </div>
      <h1 className="mb-3 text-2xl font-semibold">Что-то пошло не так</h1>
      <p className="mb-8 text-stone-600">
        {error.message || "Попробуйте обновить страницу или вернуться на главную."}
      </p>
      <div className="flex flex-wrap justify-center gap-3">
        <button type="button" onClick={reset} className="btn-primary">
          Попробовать снова
        </button>
        <Link href="/" className="btn-secondary">
          На главную
        </Link>
      </div>
    </div>
  );
}
