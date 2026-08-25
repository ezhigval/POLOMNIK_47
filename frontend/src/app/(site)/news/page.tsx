import { Suspense } from "react";
import type { Metadata } from "next";
import { NewsFeed } from "@/components/news-feed";
import { PageIntro } from "@/components/page-intro";
import { listPublicNews } from "@/lib/api/news";
import { toFeedArticle } from "@/lib/news";
import { siteConfig } from "@/lib/site-config";

export const metadata: Metadata = {
  title: "Новости",
  description: `Новости и статьи паломнической службы «${siteConfig.name}».`,
  alternates: { canonical: "/news" },
  openGraph: {
    title: "Новости",
    description: `Новости и статьи паломнической службы «${siteConfig.name}».`,
    url: "/news",
  },
};

export default async function NewsPage() {
  let articles: ReturnType<typeof toFeedArticle>[] = [];
  try {
    const published = await listPublicNews();
    articles = published.map(toFeedArticle);
  } catch {
    articles = [];
  }

  return (
    <div className="mx-auto max-w-6xl space-y-8 px-4 py-8 sm:py-10">
      <PageIntro
        title="Новости"
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
    <div className="space-y-5" aria-hidden="true">
      <div className="grid gap-5 lg:grid-cols-3">
        <div className="h-80 animate-pulse rounded-3xl bg-stone-200 lg:col-span-2" />
        <div className="grid gap-5">
          <div className="h-36 animate-pulse rounded-3xl bg-stone-200" />
          <div className="h-36 animate-pulse rounded-3xl bg-stone-200" />
        </div>
      </div>
    </div>
  );
}
