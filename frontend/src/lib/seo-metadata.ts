import type { Metadata } from "next";

type PageMetaInput = {
  title: string;
  description: string;
  canonical: string;
  ogType?: "website" | "article";
  publishedTime?: string;
  images?: string[];
};

/** Shared Open Graph + Twitter from existing page title/description (no generated copy). */
export function buildPublicPageMetadata(input: PageMetaInput): Metadata {
  const images = input.images?.filter(Boolean).map((url) => ({ url }));
  return {
    title: input.title,
    description: input.description,
    alternates: { canonical: input.canonical },
    openGraph: {
      title: input.title,
      description: input.description,
      url: input.canonical,
      type: input.ogType ?? "website",
      ...(input.publishedTime ? { publishedTime: input.publishedTime } : {}),
      ...(images?.length ? { images } : {}),
    },
    twitter: {
      card: images?.length ? "summary_large_image" : "summary",
      title: input.title,
      description: input.description,
      ...(images?.[0] ? { images: [images[0].url] } : {}),
    },
  };
}
