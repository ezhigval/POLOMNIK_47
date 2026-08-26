"use client";

import { FormEvent, useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

type Comment = {
  id: string;
  author: string;
  body: string;
  created_at: string;
};

type NewsCommentsProps = {
  slug: string;
  loggedIn: boolean;
};

export function NewsComments({ slug, loggedIn }: NewsCommentsProps) {
  const pathname = usePathname();
  const [comments, setComments] = useState<Comment[]>([]);
  const [body, setBody] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(`/api/news/${encodeURIComponent(slug)}/comments`)
      .then((response) => response.json())
      .then((payload) => setComments(payload.data ?? []))
      .catch(() => setComments([]));
  }, [slug]);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    if (!loggedIn || !body.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/news/${encodeURIComponent(slug)}/comments`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ body }),
      });
      if (response.status === 401) {
        setError("Войдите, чтобы комментировать");
        return;
      }
      if (!response.ok) {
        setError("Не удалось отправить комментарий");
        return;
      }
      const payload = await response.json();
      setComments((prev) => [...prev, payload.data]);
      setBody("");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-4 border-t border-stone-100 pt-4">
      <h3 className="text-sm font-semibold text-stone-900">Комментарии</h3>
      {comments.length === 0 ? (
        <p className="text-sm text-stone-500">Пока нет комментариев.</p>
      ) : (
        <ul className="space-y-3">
          {comments.map((comment) => (
            <li key={comment.id} className="rounded-xl bg-stone-50 p-3 text-sm">
              <p className="font-medium text-stone-900">{comment.author}</p>
              <p className="mt-1 leading-6 text-stone-700">{comment.body}</p>
            </li>
          ))}
        </ul>
      )}
      {loggedIn ? (
        <form onSubmit={onSubmit} className="space-y-2">
          <textarea
            value={body}
            onChange={(event) => setBody(event.target.value)}
            rows={3}
            className="input-field w-full text-sm"
            placeholder="Ваш комментарий"
          />
          {error ? <p className="text-sm text-red-600">{error}</p> : null}
          <button type="submit" disabled={loading || !body.trim()} className="btn-primary px-4 py-2 text-sm">
            {loading ? "Отправляем..." : "Комментировать"}
          </button>
        </form>
      ) : (
        <p className="text-sm text-stone-600">
          <Link
            href={`/account/login?returnUrl=${encodeURIComponent(pathname)}`}
            className="font-medium text-brand-800 hover:text-brand-900"
          >
            Войдите
          </Link>
          , чтобы оставить комментарий.
        </p>
      )}
    </div>
  );
}
