import { Suspense } from "react";
import { NewsInfiniteFeed } from "@/components/news-infinite-feed";
import { NewsCollectionStructuredData } from "@/components/structured-data";
import { PageIntro } from "@/components/page-intro";
import { listPublicNews } from "@/lib/api/news";
import { toFeedArticle } from "@/lib/news";
import { buildPublicPageMetadata } from "@/lib/seo-metadata";
import { siteConfig } from "@/lib/site-config";
import { getSessionUser } from "@/lib/auth/session";

const description = `Новости и статьи паломнической службы «${siteConfig.name}».`;

export const metadata = buildPublicPageMetadata({
  title: "Новости",
  description,
  canonical: "/news",
});

export default async function NewsPage() {
  const [sessionUser, published] = await Promise.all([
    getSessionUser(),
    listPublicNews().catch(() => []),
  ]);
  const articles = published.map(toFeedArticle);

  return (
    <div className="mx-auto max-w-3xl space-y-8 px-4 py-8 sm:py-10">
      <NewsCollectionStructuredData articles={articles} />
      <PageIntro
        title="Новости"
        description="Лента событий службы — лайки без регистрации, комментарии после входа в кабинет."
      />

      {articles.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-stone-300 bg-white p-8 text-center text-stone-500">
          Новостей пока нет.
        </div>
      ) : (
        <Suspense fallback={<FeedSkeleton />}>
          <NewsInfiniteFeed
            key={articles.map((article) => article.slug).join("|")}
            articles={articles}
            loggedIn={Boolean(sessionUser)}
          />
        </Suspense>
      )}
    </div>
  );
}

function FeedSkeleton() {
  return (
    <div className="space-y-6" aria-hidden="true">
      <div className="h-64 animate-pulse rounded-3xl bg-stone-200" />
      <div className="h-64 animate-pulse rounded-3xl bg-stone-200" />
    </div>
  );
}
