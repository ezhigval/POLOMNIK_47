"use client";

import { useState, useTransition } from "react";
import { toggleNewsPinAction } from "@/app/management/actions";

type NewsPinStarButtonProps = {
  id: string;
  isPinned: boolean;
  onError?: (message: string | null) => void;
};

export function NewsPinStarButton({ id, isPinned, onError }: NewsPinStarButtonProps) {
  const [pending, startTransition] = useTransition();
  const [localError, setLocalError] = useState<string | null>(null);

  function onClick() {
    setLocalError(null);
    onError?.(null);
    startTransition(async () => {
      try {
        await toggleNewsPinAction(id, !isPinned);
      } catch (err) {
        const message = err instanceof Error ? err.message : "Не удалось изменить закрепление";
        setLocalError(message);
        onError?.(message);
      }
    });
  }

  const label = isPinned ? "Открепить новость" : "Закрепить новость";

  return (
    <div className="flex flex-col items-start gap-1">
      <button
        type="button"
        onClick={onClick}
        disabled={pending}
        aria-pressed={isPinned}
        aria-label={label}
        title={label}
        className={`rounded-full p-1 transition hover:bg-brand-50 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-700 disabled:opacity-60 ${
          isPinned ? "text-brand-800" : "text-stone-400 hover:text-brand-700"
        }`}
      >
        <StarIcon filled={isPinned} />
      </button>
      {localError ? <span className="max-w-[11rem] text-xs text-red-700">{localError}</span> : null}
    </div>
  );
}

function StarIcon({ filled }: { filled: boolean }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      className="size-6"
      aria-hidden="true"
      fill={filled ? "currentColor" : "none"}
      stroke="currentColor"
      strokeWidth={1.6}
      strokeLinejoin="round"
    >
      <path d="M12 3.6 14.47 9.1l6.03.54-4.56 3.92 1.4 5.84L12 16.4l-5.34 3 1.4-5.84-4.56-3.92 6.03-.54L12 3.6Z" />
    </svg>
  );
}
