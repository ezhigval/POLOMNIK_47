export const siteConfig = {
  name: process.env.NEXT_PUBLIC_SITE_NAME ?? "POLOMNIK 47",
  tagline: process.env.NEXT_PUBLIC_SITE_TAGLINE ?? "Паломническая служба",
  url: process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000",
  description:
    process.env.NEXT_PUBLIC_SITE_DESCRIPTION ??
    "Паломнические поездки в монастыри и святые места России. Сопровождение духовника, комфортный транспорт, прозрачная стоимость.",
  region: process.env.NEXT_PUBLIC_SITE_REGION ?? "RU",
  departureCity: process.env.NEXT_PUBLIC_DEPARTURE_CITY ?? "Санкт-Петербург",
};

export function absoluteUrl(path: string): string {
  const base = siteConfig.url.replace(/\/$/, "");
  const normalized = path.startsWith("/") ? path : `/${path}`;
  return `${base}${normalized}`;
}
