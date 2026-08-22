import type { Tour } from "@/lib/api/tours";

const slugCovers: Record<string, string> = {
  "optina-pustyn":
    "https://images.unsplash.com/photo-1548013146-724f68d1ddac?auto=format&fit=crop&w=1200&q=80",
  diveevo:
    "https://images.unsplash.com/photo-1605647540924-852290f6b0d5?auto=format&fit=crop&w=1200&q=80",
};

const locationKeywords: { match: RegExp; url: string }[] = [
  { match: /оптин|калуж/i, url: slugCovers["optina-pustyn"] },
  { match: /дивеев|нижегород/i, url: slugCovers.diveevo },
  { match: /солов|архангел/i, url: "https://images.unsplash.com/photo-1516026672322-bc52d61a55d5?auto=format&fit=crop&w=1200&q=80" },
  { match: /валда|тихвин/i, url: "https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?auto=format&fit=crop&w=1200&q=80" },
];

export const heroBackgroundImages = [
  "/images/hero/tikhvin-monastery.webp",
  "/images/hero/tikhvin-sunset.webp",
  "/images/hero/wooden-chapel.webp",
  "/images/hero/monastery-path.webp",
] as const;

export function getTourCoverUrl(tour: Pick<Tour, "slug" | "location" | "title" | "images">): string | null {
  if (tour.images?.[0]) {
    return tour.images[0];
  }

  if (tour.slug && slugCovers[tour.slug]) {
    return slugCovers[tour.slug];
  }

  const haystack = `${tour.location} ${tour.title}`;
  for (const entry of locationKeywords) {
    if (entry.match.test(haystack)) {
      return entry.url;
    }
  }

  return null;
}
