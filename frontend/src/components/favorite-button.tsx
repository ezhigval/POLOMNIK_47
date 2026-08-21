"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

type FavoriteButtonProps = {
  tourId: string;
  compact?: boolean;
};

export function FavoriteButton({ tourId, compact = false }: FavoriteButtonProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [saved, setSaved] = useState(false);
  const [loading, setLoading] = useState(false);
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    let cancelled = false;
    fetch("/api/favorites")
      .then(async (response) => {
        if (response.status === 401) {
          return null;
        }
        if (!response.ok) {
          return null;
        }
        const payload = await response.json();
        return payload.data as string[];
      })
      .then((ids) => {
        if (cancelled || !ids) {
          return;
        }
        setSaved(ids.includes(tourId));
        setChecked(true);
      });
    return () => {
      cancelled = true;
    };
  }, [tourId]);

  async function toggleFavorite(event: React.MouseEvent) {
    event.preventDefault();
    event.stopPropagation();
    setLoading(true);

    const response = await fetch(`/api/favorites/${tourId}`, {
      method: saved ? "DELETE" : "POST",
    });

    setLoading(false);
    if (response.status === 401) {
      const returnUrl = encodeURIComponent(pathname);
      router.push(`/account/login?returnUrl=${returnUrl}`);
      return;
    }
    if (!response.ok) {
      return;
    }

    setSaved((value) => !value);
    setChecked(true);
    router.refresh();
  }

  if (!checked) {
    return (
      <span
        className={`inline-flex items-center justify-center rounded-full border border-white/80 bg-white/90 ${
          compact ? "size-9" : "size-10"
        }`}
        aria-hidden
      />
    );
  }

  return (
    <button
      type="button"
      onClick={toggleFavorite}
      disabled={loading}
      aria-pressed={saved}
      aria-label={saved ? "Убрать из избранного" : "Добавить в избранное"}
      className={`inline-flex items-center justify-center rounded-full border transition ${
        compact ? "size-9" : "size-10"
      } ${
        saved
          ? "border-red-200 bg-red-50 text-red-600 hover:bg-red-100"
          : "border-white/80 bg-white/90 text-stone-600 hover:bg-white hover:text-red-600"
      } disabled:opacity-60`}
    >
      {saved ? "♥" : "♡"}
    </button>
  );
}

export function FavoriteLoginHint() {
  const pathname = usePathname();
  const loginHref = `/account/login?returnUrl=${encodeURIComponent(pathname)}`;

  return (
    <p className="text-xs text-stone-500">
      <Link href={loginHref} className="font-medium text-brand-800 hover:underline">
        Войдите
      </Link>
      , чтобы сохранять туры в избранное
    </p>
  );
}
