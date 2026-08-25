import Link from "next/link";
import { formatNewsDate, splitNewsLayout, type NewsArticle } from "@/lib/news";

type NewsFeedProps = {
  articles: NewsArticle[];
};

export function NewsFeed({ articles }: NewsFeedProps) {
  const { featured, side, rest } = splitNewsLayout(articles);

  return (
    <div className="space-y-8">
      {featured ? (
        <div className="grid gap-5 lg:grid-cols-3">
          <div className={side.length > 0 ? "lg:col-span-2" : "lg:col-span-3"}>
            <NewsCard article={featured} featured />
          </div>
          {side.length > 0 ? (
            <ul className="grid gap-5">
              {side.map((article) => (
                <li key={article.slug}>
                  <NewsCard article={article} />
                </li>
              ))}
            </ul>
          ) : null}
        </div>
      ) : null}

      {rest.length > 0 ? (
        <ul className="grid gap-5 sm:grid-cols-2">
          {rest.map((article) => (
            <li key={article.slug}>
              <NewsCard article={article} />
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}

function NewsCard({ article, featured = false }: { article: NewsArticle; featured?: boolean }) {
  return (
    <Link
      href={`/news/${article.slug}`}
      className="group flex h-full w-full flex-col overflow-hidden rounded-3xl border border-stone-200 bg-white text-left shadow-sm transition hover:-translate-y-0.5 hover:border-brand-200 hover:shadow-md"
    >
      <span className={`relative block overflow-hidden bg-stone-100 ${featured ? "aspect-[16/9]" : "aspect-[16/9]"}`}>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={article.image}
          alt={article.title}
          className="size-full object-cover transition duration-500 group-hover:scale-[1.03]"
        />
      </span>
      <span className={`flex flex-1 flex-col ${featured ? "p-5 sm:p-6" : "p-4 sm:p-5"}`}>
        <span className="flex flex-wrap items-center gap-2">
          <time dateTime={article.date} className="text-xs font-medium uppercase tracking-wide text-brand-800">
            {formatNewsDate(article.date)}
          </time>
          {article.pinned ? (
            <span className="rounded-full bg-brand-50 px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide text-brand-800">
              Закреплено
            </span>
          ) : null}
        </span>
        <span
          className={`mt-2 font-display font-semibold text-stone-900 group-hover:text-brand-800 ${
            featured ? "text-xl sm:text-2xl" : "text-lg"
          }`}
        >
          {article.title}
        </span>
        <span className={`mt-2 text-sm leading-6 text-stone-600 ${featured ? "line-clamp-4" : "line-clamp-3"}`}>
          {article.excerpt}
        </span>
        <span className="mt-4 text-sm font-medium text-brand-800">Читать статью →</span>
      </span>
    </Link>
  );
}
