import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { formatNewsDate, NEWS_AI_DISCLAIMER, paragraphsFromBody, toFeedArticle } from "@/lib/news";
import { getPublicNewsBySlug } from "@/lib/api/news";
import { ApiError } from "@/lib/api/client";
import { siteConfig } from "@/lib/site-config";

type PageProps = {
  params: Promise<{ slug: string }>;
};

async function loadArticle(slug: string) {
  try {
    return toFeedArticle(await getPublicNewsBySlug(slug));
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      return null;
    }
    throw error;
  }
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;
  const article = await loadArticle(slug);
  if (!article) {
    return { title: "Новость" };
  }
  return {
    title: article.title,
    description: article.excerpt,
    alternates: { canonical: `/news/${article.slug}` },
    openGraph: {
      title: article.title,
      description: article.excerpt,
      url: `/news/${article.slug}`,
    },
  };
}

export default async function NewsArticlePage({ params }: PageProps) {
  const { slug } = await params;
  const article = await loadArticle(slug);
  if (!article) {
    notFound();
  }

  const paragraphs = article.paragraphs;
  const hasDisclaimer = paragraphs.at(-1)?.trim() === NEWS_AI_DISCLAIMER;
  const body = hasDisclaimer ? paragraphs.slice(0, -1) : paragraphs;

  return (
    <article className="mx-auto max-w-3xl space-y-6 px-4 py-8 sm:py-10">
      <Link href="/news" className="text-sm text-stone-500 hover:text-brand-800">
        ← Все новости
      </Link>
      <div className="overflow-hidden rounded-3xl border border-stone-200 bg-white shadow-sm">
        {article.image ? (
          <div className="aspect-[16/8] bg-stone-100">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src={article.image} alt={article.title} className="size-full object-cover" />
          </div>
        ) : null}
        <div className="space-y-4 px-5 py-6 sm:px-8 sm:py-8">
          <time dateTime={article.date} className="text-xs font-medium uppercase tracking-wide text-brand-800">
            {formatNewsDate(article.date)}
          </time>
          <h1 className="font-display text-3xl font-semibold text-stone-900">{article.title}</h1>
          <p className="sr-only">{siteConfig.name}</p>
          <div className="space-y-4 text-sm leading-7 text-stone-700 sm:text-base">
            {(body.length > 0 ? body : paragraphsFromBody(article.excerpt)).map((paragraph) => (
              <p key={paragraph}>{paragraph}</p>
            ))}
          </div>
          {hasDisclaimer ? (
            <p className="border-t border-stone-100 pt-4 text-xs leading-5 text-stone-400">{NEWS_AI_DISCLAIMER}</p>
          ) : null}
        </div>
      </div>
    </article>
  );
}
