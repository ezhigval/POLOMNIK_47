"use client";

import { useEffect, useRef, useState } from "react";
import type { NewsArticle } from "@/lib/news";
import { NewsFeedItem } from "@/components/news-feed-item";

const PAGE_SIZE = 5;

type NewsInfiniteFeedProps = {
  articles: NewsArticle[];
  loggedIn: boolean;
};

export function NewsInfiniteFeed({ articles, loggedIn }: NewsInfiniteFeedProps) {
  const [visibleCount, setVisibleCount] = useState(PAGE_SIZE);
  const sentinelRef = useRef<HTMLDivElement>(null);
  const articlesKey = articles.map((article) => article.slug).join("|");

  useEffect(() => {
    const element = sentinelRef.current;
    if (!element || visibleCount >= articles.length) {
      return;
    }

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setVisibleCount((current) => Math.min(current + PAGE_SIZE, articles.length));
        }
      },
      { rootMargin: "240px" },
    );
    observer.observe(element);
    return () => observer.disconnect();
  }, [articles.length, articlesKey, visibleCount]);

  const visible = articles.slice(0, visibleCount);

  return (
    <div className="space-y-6">
      {visible.map((article) => (
        <NewsFeedItem key={article.slug} article={article} loggedIn={loggedIn} />
      ))}
      {visibleCount < articles.length ? (
        <div ref={sentinelRef} className="h-10 animate-pulse rounded-2xl bg-stone-100" aria-hidden />
      ) : null}
    </div>
  );
}
