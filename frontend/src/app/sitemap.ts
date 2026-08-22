import type { MetadataRoute } from "next";
import { getTours } from "@/lib/api/tours";
import { absoluteUrl } from "@/lib/site-config";

export const dynamic = "force-dynamic";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const staticPages: MetadataRoute.Sitemap = [
    { url: absoluteUrl("/"), changeFrequency: "daily", priority: 1 },
    { url: absoluteUrl("/search"), changeFrequency: "daily", priority: 0.9 },
    { url: absoluteUrl("/reviews"), changeFrequency: "weekly", priority: 0.7 },
    { url: absoluteUrl("/news"), changeFrequency: "weekly", priority: 0.7 },
    { url: absoluteUrl("/support"), changeFrequency: "monthly", priority: 0.5 },
    { url: absoluteUrl("/privacy"), changeFrequency: "monthly", priority: 0.3 },
  ];

  try {
    if (!process.env.API_INTERNAL_URL) {
      return staticPages;
    }
    const response = await getTours({ limit: "100" });
    const tourPages = response.data.map((tour) => ({
      url: absoluteUrl(`/tours/${tour.id}`),
      changeFrequency: "weekly" as const,
      priority: 0.8,
    }));
    return [...staticPages, ...tourPages];
  } catch {
    return staticPages;
  }
}
