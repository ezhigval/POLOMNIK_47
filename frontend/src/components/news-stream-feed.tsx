"use client";

import Link from "next/link";
import { useState } from "react";
import { NewsComments } from "@/components/news-comments";
import { NewsLikeButton } from "@/components/news-like-button";
import { PhotoLightbox } from "@/components/photo-lightbox";
import { linkifyText } from "@/lib/linkify";
import { formatNewsDate, newsPreviewImageSrc, type NewsArticle } from "@/lib/news";

type NewsStreamFeedProps = {
  articles: NewsArticle[];
  loggedIn: boolean;
};

export function NewsStreamFeed({ articles, loggedIn }: NewsStreamFeedProps) {
  return (
    <div className="mx-auto max-w-2xl space-y-8">
      {articles.map((article) => (
        <NewsStreamCard key={article.slug} article={article} loggedIn={loggedIn} />
      ))}
    </div>
  );
}

function NewsStreamCard({ article, loggedIn }: { article: NewsArticle; loggedIn: boolean }) {
  const [lightboxSrc, setLightboxSrc] = useState<string | null>(null);
  const previewImage = newsPreviewImageSrc(article);
  const images = previewImage ? [previewImage] : [];

  return (
    <article className="overflow-hidden rounded-2xl border border-stone-200 bg-white shadow-sm">
      <header className="border-b border-stone-100 px-5 py-4">
        <div className="flex flex-wrap items-center gap-2 text-xs font-medium uppercase tracking-wide text-brand-800">
          <time dateTime={article.date}>{formatNewsDate(article.date)}</time>
          {article.pinned ? (
            <span className="rounded-full bg-brand-50 px-2 py-0.5 text-[10px] text-brand-800">Закреплено</span>
          ) : null}
        </div>
        <h2 className="mt-2 font-display text-xl font-semibold text-stone-900">
          <Link href={`/news/${article.slug}`} className="hover:text-brand-800">
            {article.title}
          </Link>
        </h2>
        {article.excerpt ? <p className="mt-2 text-sm text-stone-600">{article.excerpt}</p> : null}
      </header>

      {images.length > 0 ? (
        <div className="space-y-1">
          {images.map((src) => (
            <button
              key={src}
              type="button"
              className="block w-full cursor-zoom-in"
              onClick={() => setLightboxSrc(src)}
            >
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={src} alt="" className="block w-full" />
            </button>
          ))}
        </div>
      ) : null}

      {article.paragraphs.length > 0 ? (
        <div className="space-y-3 px-5 py-4 text-sm">
          {article.paragraphs.slice(0, 4).map((paragraph) => (
            <p key={paragraph} className="leading-7 text-stone-700">
              {linkifyText(paragraph)}
            </p>
          ))}
          {article.paragraphs.length > 4 ? (
            <Link href={`/news/${article.slug}`} className="text-sm font-medium text-brand-800 hover:text-brand-900">
              Читать полностью →
            </Link>
          ) : null}
        </div>
      ) : null}

      <footer className="space-y-4 px-5 pb-5">
        <NewsLikeButton slug={article.slug} />
        <NewsComments slug={article.slug} loggedIn={loggedIn} />
      </footer>

      <PhotoLightbox
        src={lightboxSrc ?? ""}
        alt={article.title}
        open={lightboxSrc != null}
        onClose={() => setLightboxSrc(null)}
      />
    </article>
  );
}
