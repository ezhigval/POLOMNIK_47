import Link from "next/link";
import { NewsCompactList } from "@/components/news-compact-list";
import { SectionHeading } from "@/components/section-heading";
import { listPublicNews } from "@/lib/api/news";
import { toFeedArticle } from "@/lib/news";

const homepageNewsLimit = 6;

export async function HomeNewsSection() {
  let articles: ReturnType<typeof toFeedArticle>[] = [];
  try {
    const published = await listPublicNews();
    articles = published.map(toFeedArticle).slice(0, homepageNewsLimit);
  } catch {
    return null;
  }

  if (articles.length === 0) {
    return null;
  }

  return (
    <section className="scroll-mt-24">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <SectionHeading
          eyebrow="Новости"
          title="События службы"
          description="События службы, маршруты и святыйни."
        />
        <Link href="/news" className="text-sm font-medium text-brand-800 hover:text-brand-900">
          Все новости →
        </Link>
      </div>
      <div className="mt-6">
        <NewsCompactList articles={articles} />
      </div>
    </section>
  );
}

export function HomeNewsSkeleton() {
  return (
    <div className="space-y-4" aria-busy="true">
      <div className="h-10 w-64 animate-pulse rounded-lg bg-stone-200" />
      <div className="h-48 animate-pulse rounded-2xl bg-stone-200" />
    </div>
  );
}
