import Link from "next/link";
import { formatNewsDate, type NewsArticle } from "@/lib/news";

type NewsFeedProps = {
  articles: NewsArticle[];
};

export function NewsFeed({ articles }: NewsFeedProps) {
  return (
    <ul className="grid gap-5 sm:grid-cols-2">
      {articles.map((article, index) => (
        <li key={article.slug} className={index === 0 ? "sm:col-span-2" : undefined}>
          <Link
            href={`/news/${article.slug}`}
            className="group flex h-full w-full flex-col overflow-hidden rounded-3xl border border-stone-200 bg-white text-left shadow-sm transition hover:-translate-y-0.5 hover:border-brand-200 hover:shadow-md"
          >
            <span className="relative block aspect-[16/9] overflow-hidden bg-stone-100">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={article.image}
                alt={article.title}
                className="size-full object-cover transition duration-500 group-hover:scale-[1.03]"
              />
            </span>
            <span className="flex flex-1 flex-col p-5 sm:p-6">
              <time dateTime={article.date} className="text-xs font-medium uppercase tracking-wide text-brand-800">
                {formatNewsDate(article.date)}
              </time>
              <span className="mt-2 font-display text-xl font-semibold text-stone-900 group-hover:text-brand-800 sm:text-2xl">
                {article.title}
              </span>
              {article.photoStrip.length === 0 ? (
                <span className="mt-2 line-clamp-3 text-sm leading-6 text-stone-600">{article.excerpt}</span>
              ) : null}
              <span className="mt-4 text-sm font-medium text-brand-800">Читать статью →</span>
            </span>
          </Link>
        </li>
      ))}
    </ul>
  );
}
