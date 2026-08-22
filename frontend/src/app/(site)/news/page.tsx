import { Suspense } from "react";
import type { Metadata } from "next";
import { NewsFeed } from "@/components/news-feed";
import { PageIntro } from "@/components/page-intro";
import { listPublicNews } from "@/lib/api/news";
import { newsArticles, toFeedArticle } from "@/lib/news";
import { siteConfig } from "@/lib/site-config";

export const metadata: Metadata = {
  title: "Новостная лента",
  description: `Новости и статьи паломнической службы «${siteConfig.name}».`,
  alternates: { canonical: "/news" },
};

export default async function NewsPage() {
  let articles = newsArticles;
  try {
    const published = await listPublicNews();
    articles = published.map(toFeedArticle);
  } catch {
    articles = newsArticles;
  }

  return (
    <div className="mx-auto max-w-6xl space-y-8 px-4 py-8 sm:py-10">
      <PageIntro
        title="Новостная лента"
        description="События службы, маршруты и святыни — откройте карточку, чтобы прочитать статью."
      />

      {articles.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-stone-300 bg-white p-8 text-center text-stone-500">
          Новостей пока нет.
        </div>
      ) : (
        <Suspense fallback={<NewsFeedSkeleton />}>
          <NewsFeed articles={articles} />
        </Suspense>
      )}
    </div>
  );
}

function NewsFeedSkeleton() {
  return (
    <div className="grid gap-5 sm:grid-cols-2" aria-hidden="true">
      {Array.from({ length: 4 }).map((_, index) => (
        <div
          key={index}
          className={`h-72 animate-pulse rounded-3xl bg-stone-200 ${index === 0 ? "sm:col-span-2" : ""}`}
        />
      ))}
    </div>
  );
}
