"use client";

import Link from "next/link";
import { formatNewsDate, isNewsImageSrc, NEWS_AI_DISCLAIMER, type NewsArticle } from "@/lib/news";
import { linkifyText } from "@/lib/linkify";
import { NewsLikeButton } from "@/components/news-like-button";
import { NewsComments } from "@/components/news-comments";
import { NewsClickableImage } from "@/components/news-clickable-image";

type NewsFeedItemProps = {
  article: NewsArticle;
  loggedIn: boolean;
};

export function NewsFeedItem({ article, loggedIn }: NewsFeedItemProps) {
  const hasPhotoStrip = article.photoStrip.length > 0;
  const paragraphs = article.paragraphs;
  const hasDisclaimer = paragraphs.at(-1)?.trim() === NEWS_AI_DISCLAIMER;
  const body = hasDisclaimer ? paragraphs.slice(0, -1) : paragraphs;
  const previewParagraphs = body.slice(0, 3);

  return (
    <article className="overflow-hidden rounded-3xl border border-stone-200 bg-white shadow-sm">
      <div className="space-y-4 px-5 py-5 sm:px-6 sm:py-6">
        <header className="space-y-2">
          <div className="flex flex-wrap items-center gap-2">
            <time dateTime={article.date} className="text-xs font-medium uppercase tracking-wide text-brand-800">
              {formatNewsDate(article.date)}
            </time>
            {article.pinned ? (
              <span className="rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide text-amber-900">
                Закреплено
              </span>
            ) : null}
          </div>
          <h2 className="font-display text-xl font-semibold text-stone-900 sm:text-2xl">
            <Link href={`/news/${article.slug}`} className="hover:text-brand-800">
              {article.title}
            </Link>
          </h2>
        </header>

        {hasPhotoStrip ? (
          <div className="space-y-2">
            {article.photoStrip.map((src) => (
              <NewsClickableImage key={src} src={src} alt="" className="block h-auto w-full rounded-2xl" />
            ))}
          </div>
        ) : article.image ? (
          <NewsClickableImage src={article.image} alt={article.title} className="aspect-[16/9] w-full rounded-2xl object-cover" />
        ) : null}

        {!hasPhotoStrip && previewParagraphs.length > 0 ? (
          <div className="space-y-3 text-sm sm:text-base">
            {previewParagraphs.map((paragraph) =>
              isNewsImageSrc(paragraph) ? (
                <NewsClickableImage key={paragraph} src={paragraph.trim()} alt="" className="mx-auto h-auto w-full rounded-2xl object-contain" />
              ) : (
                <p key={paragraph} className="leading-7 text-stone-700">
                  {linkifyText(paragraph)}
                </p>
              ),
            )}
          </div>
        ) : null}

        {!hasPhotoStrip && article.excerpt ? (
          <p className="text-sm leading-6 text-stone-600">{article.excerpt}</p>
        ) : null}

        <div className="flex flex-wrap items-center gap-3">
          <NewsLikeButton slug={article.slug} />
          <Link href={`/news/${article.slug}`} className="text-sm font-medium text-brand-800 hover:text-brand-900">
            Читать полностью →
          </Link>
        </div>

        <NewsComments slug={article.slug} loggedIn={loggedIn} />
      </div>
    </article>
  );
}
