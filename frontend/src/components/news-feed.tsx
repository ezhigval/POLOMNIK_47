"use client";

import { useCallback, useEffect } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { formatNewsDate, NEWS_AI_DISCLAIMER, type NewsArticle } from "@/lib/news";

type NewsFeedProps = {
  articles: NewsArticle[];
};

export function NewsFeed({ articles }: NewsFeedProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const openSlug = searchParams.get("article");
  const openArticle = articles.find((article) => article.slug === openSlug) ?? null;

  const closeArticlePopup = useCallback(() => {
    const params = new URLSearchParams(searchParams.toString());
    params.delete("article");
    const query = params.toString();
    router.replace(query ? `${pathname}?${query}` : pathname, { scroll: false });
  }, [pathname, router, searchParams]);

  function openArticlePopup(slug: string) {
    const params = new URLSearchParams(searchParams.toString());
    params.set("article", slug);
    router.replace(`${pathname}?${params.toString()}`, { scroll: false });
  }

  useEffect(() => {
    if (!openArticle) {
      return;
    }

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        closeArticlePopup();
      }
    }

    window.addEventListener("keydown", onKeyDown);
    return () => {
      document.body.style.overflow = previousOverflow;
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [closeArticlePopup, openArticle]);

  return (
    <>
      <ul className="grid gap-5 sm:grid-cols-2">
        {articles.map((article, index) => (
          <li key={article.slug} className={index === 0 ? "sm:col-span-2" : undefined}>
            <button
              type="button"
              onClick={() => openArticlePopup(article.slug)}
              className="group flex h-full w-full flex-col overflow-hidden rounded-3xl border border-stone-200 bg-white text-left shadow-sm transition hover:-translate-y-0.5 hover:border-brand-200 hover:shadow-md"
            >
              <span className="relative block aspect-[16/9] overflow-hidden bg-stone-100">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={article.image}
                  alt={article.title}
                  className="size-full object-cover transition duration-500 group-hover:scale-[1.03]"
                />
              </span>
              <span className="flex flex-1 flex-col p-5 sm:p-6">
                <time dateTime={article.date} className="text-xs font-medium uppercase tracking-wide text-brand-800">
                  {formatNewsDate(article.date)}
                </time>
                <span className="mt-2 font-display text-xl font-semibold text-stone-900 group-hover:text-brand-800 sm:text-2xl">
                  {article.title}
                </span>
                <span className="mt-2 line-clamp-3 text-sm leading-6 text-stone-600">{article.excerpt}</span>
                <span className="mt-4 text-sm font-medium text-brand-800">Читать статью →</span>
              </span>
            </button>
          </li>
        ))}
      </ul>

      {openArticle ? (
        <div className="fixed inset-0 z-50 flex items-end justify-center p-0 sm:items-center sm:p-6">
          <button
            type="button"
            className="absolute inset-0 bg-stone-950/60 backdrop-blur-[2px]"
            aria-label="Закрыть статью"
            onClick={closeArticlePopup}
          />
          <article
            role="dialog"
            aria-modal="true"
            aria-labelledby="news-article-title"
            className="relative z-10 flex max-h-[92vh] w-full max-w-2xl flex-col overflow-hidden rounded-t-3xl bg-white shadow-2xl sm:rounded-3xl"
          >
            <div className="relative aspect-[16/8] shrink-0 bg-stone-100">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={openArticle.image} alt={openArticle.title} className="size-full object-cover" />
              <button
                type="button"
                onClick={closeArticlePopup}
                className="absolute right-3 top-3 inline-flex size-10 items-center justify-center rounded-full bg-white/90 text-stone-800 shadow-sm transition hover:bg-white"
                aria-label="Закрыть"
              >
                <svg viewBox="0 0 24 24" className="size-5" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M6 6l12 12M6 18L18 6" />
                </svg>
              </button>
            </div>
            <div className="overflow-y-auto px-5 py-6 sm:px-8 sm:py-8">
              <time dateTime={openArticle.date} className="text-xs font-medium uppercase tracking-wide text-brand-800">
                {formatNewsDate(openArticle.date)}
              </time>
              <h2 id="news-article-title" className="mt-2 font-display text-2xl font-semibold text-stone-900 sm:text-3xl">
                {openArticle.title}
              </h2>
              <div className="mt-5 space-y-4 text-sm leading-7 text-stone-700 sm:text-base">
                {articleParagraphs(openArticle.paragraphs).map((paragraph) => (
                  <p key={paragraph}>{paragraph}</p>
                ))}
              </div>
              {hasAiDisclaimer(openArticle.paragraphs) ? (
                <p className="mt-6 border-t border-stone-100 pt-4 text-xs leading-5 text-stone-400">
                  {NEWS_AI_DISCLAIMER}
                </p>
              ) : null}
            </div>
          </article>
        </div>
      ) : null}
    </>
  );
}

function articleParagraphs(paragraphs: string[]): string[] {
  if (hasAiDisclaimer(paragraphs)) {
    return paragraphs.slice(0, -1);
  }
  return paragraphs;
}

function hasAiDisclaimer(paragraphs: string[]): boolean {
  return paragraphs.at(-1)?.trim() === NEWS_AI_DISCLAIMER;
}
