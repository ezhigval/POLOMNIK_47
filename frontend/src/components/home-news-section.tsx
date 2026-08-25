import Link from "next/link";
import { NewsFeed } from "@/components/news-feed";
import { SectionHeading } from "@/components/section-heading";
import { listPublicNews } from "@/lib/api/news";
import { splitNewsLayout, toFeedArticle } from "@/lib/news";

const homepageRestLimit = 4;

export async function HomeNewsSection() {
  let articles: ReturnType<typeof toFeedArticle>[] = [];
  try {
    const published = await listPublicNews();
    articles = published.map(toFeedArticle);
  } catch {
    return null;
  }

  if (articles.length === 0) {
    return null;
  }

  const { featured, side, rest } = splitNewsLayout(articles);
  const homepageArticles = [featured, ...side, ...rest.slice(0, homepageRestLimit)].filter(
    (item): item is NonNullable<typeof item> => item != null,
  );
  const hasMore = rest.length > homepageRestLimit;

  return (
    <section className="scroll-mt-24">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <SectionHeading
          eyebrow="Новости"
          title="События службы"
          description="События службы, маршруты и святыни."
        />
        <Link href="/news" className="text-sm font-medium text-brand-800 hover:text-brand-900">
          Все новости →
        </Link>
      </div>
      <div className="mt-8">
        <NewsFeed articles={homepageArticles} />
      </div>
      {hasMore ? (
        <p className="mt-6 text-center">
          <Link href="/news" className="text-sm font-medium text-brand-800 hover:text-brand-900">
            Смотреть все новости →
          </Link>
        </p>
      ) : null}
    </section>
  );
}

export function HomeNewsSkeleton() {
  return (
    <div className="space-y-5" aria-busy="true">
      <div className="h-10 w-64 animate-pulse rounded-lg bg-stone-200" />
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
