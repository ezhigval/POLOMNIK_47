import Link from "next/link";
import { formatNewsDate, type NewsArticle } from "@/lib/news";

type NewsCompactListProps = {
  articles: NewsArticle[];
};

export function NewsCompactList({ articles }: NewsCompactListProps) {
  return (
    <div className="overflow-hidden rounded-2xl border border-stone-200 bg-white shadow-sm">
      <ul className="divide-y divide-stone-100">
        {articles.map((article) => (
          <li key={article.slug}>
            <Link
              href={`/news/${article.slug}`}
              className="flex flex-col gap-1 px-4 py-3 transition hover:bg-stone-50 sm:flex-row sm:items-center sm:justify-between sm:gap-4"
            >
              <div className="min-w-0 flex-1">
                <p className="font-medium leading-snug text-stone-900 group-hover:text-brand-800">
                  {article.title}
                  {article.pinned ? (
                    <span className="ml-2 inline-flex rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-900">
                      Закреплено
                    </span>
                  ) : null}
                </p>
                {article.excerpt ? (
                  <p className="mt-1 line-clamp-1 text-sm text-stone-500">{article.excerpt}</p>
                ) : null}
              </div>
              <time
                dateTime={article.date}
                className="shrink-0 text-xs font-medium uppercase tracking-wide text-stone-400 sm:text-sm sm:normal-case sm:tracking-normal sm:text-stone-500"
              >
                {formatNewsDate(article.date)}
              </time>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
