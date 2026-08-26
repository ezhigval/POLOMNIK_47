import type { MetadataRoute } from "next";
import { listPublishedPages, cmsPublicPath } from "@/lib/api/cms";
import { fetchLegalDocuments } from "@/lib/api/legal";
import { listPublicNews } from "@/lib/api/news";
import { getTours } from "@/lib/api/tours";
import { legalDocumentPaths, type LegalDocumentType } from "@/lib/operator-config";
import { absoluteUrl } from "@/lib/site-config";
import { tourPath } from "@/lib/tour-path";

export const dynamic = "force-dynamic";

function entry(
  path: string,
  options: Pick<MetadataRoute.Sitemap[number], "lastModified" | "changeFrequency" | "priority">,
): MetadataRoute.Sitemap[number] {
  return {
    url: absoluteUrl(path),
    ...options,
  };
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const staticPages: MetadataRoute.Sitemap = [
    entry("/", { changeFrequency: "daily", priority: 1 }),
    entry("/search", { changeFrequency: "daily", priority: 0.9 }),
    entry("/reviews", { changeFrequency: "weekly", priority: 0.7 }),
    entry("/news", { changeFrequency: "weekly", priority: 0.7 }),
    entry("/support", { changeFrequency: "monthly", priority: 0.5 }),
    entry("/legal", { changeFrequency: "monthly", priority: 0.3 }),
  ];

  if (!process.env.API_INTERNAL_URL) {
    return staticPages;
  }

  const [tours, news, legal, cms] = await Promise.all([
    getTours({ limit: "100" }).catch(() => ({ data: [] as Awaited<ReturnType<typeof getTours>>["data"] })),
    listPublicNews().catch(() => []),
    fetchLegalDocuments().catch(() => []),
    listPublishedPages().catch(() => []),
  ]);

  const tourPages = tours.data.map((tour) =>
    entry(tourPath(tour), { changeFrequency: "weekly", priority: 0.8 }),
  );

  const newsPages = news.map((article) =>
    entry(`/news/${article.slug}`, {
      lastModified: article.published_at ? new Date(article.published_at) : undefined,
      changeFrequency: "weekly",
      priority: 0.6,
    }),
  );

  const legalPages = legal.flatMap((doc) => {
    const path = legalDocumentPaths[doc.type as LegalDocumentType];
    if (!path) {
      return [];
    }
    return [
      entry(path, {
        lastModified: doc.updated_at ? new Date(doc.updated_at) : undefined,
        changeFrequency: "yearly",
        priority: 0.2,
      }),
    ];
  });

  const cmsPages = cms.flatMap((page) => {
    const path = cmsPublicPath(page);
    if (!path) {
      return [];
    }
    return [
      entry(path, {
        lastModified: page.updated_at ? new Date(page.updated_at) : undefined,
        changeFrequency: "weekly",
        priority: 0.5,
      }),
    ];
  });

  return [...staticPages, ...tourPages, ...newsPages, ...legalPages, ...cmsPages];
}
