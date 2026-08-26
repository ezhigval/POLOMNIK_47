"use client";

import { useEffect, useState } from "react";

type NewsLikeButtonProps = {
  slug: string;
};

type LikeState = {
  like_count: number;
  liked_by_you: boolean;
};

export function NewsLikeButton({ slug }: NewsLikeButtonProps) {
  const [state, setState] = useState<LikeState | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetch(`/api/news/${encodeURIComponent(slug)}/likes`)
      .then((response) => response.json())
      .then((body) => setState(body.data))
      .catch(() => setState({ like_count: 0, liked_by_you: false }));
  }, [slug]);

  async function toggle() {
    setLoading(true);
    try {
      const response = await fetch(`/api/news/${encodeURIComponent(slug)}/likes`, { method: "POST" });
      if (response.ok) {
        const body = await response.json();
        setState(body.data);
      }
    } finally {
      setLoading(false);
    }
  }

  const count = state?.like_count ?? 0;
  const liked = state?.liked_by_you ?? false;

  return (
    <button
      type="button"
      onClick={toggle}
      disabled={loading}
      aria-pressed={liked}
      className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-sm font-medium transition ${
        liked
          ? "border-brand-200 bg-brand-50 text-brand-900"
          : "border-stone-200 bg-white text-stone-700 hover:border-brand-200"
      }`}
    >
      <span aria-hidden>{liked ? "♥" : "♡"}</span>
      <span>{count > 0 ? count : "Нравится"}</span>
    </button>
  );
}
