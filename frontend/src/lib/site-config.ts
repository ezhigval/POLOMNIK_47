export const siteConfig = {
  name: process.env.NEXT_PUBLIC_SITE_NAME ?? 'Под Покровом Божией Матери "Тихвинская"',
  fullName:
    process.env.NEXT_PUBLIC_SITE_FULL_NAME ??
    "РПЦ Тихвинская Епархия. Паломническая служба «Под покровом Божией Матери «Тихвинская»»",
  tagline: process.env.NEXT_PUBLIC_SITE_TAGLINE ?? "Паломническая служба",
  url: process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000",
  description:
    process.env.NEXT_PUBLIC_SITE_DESCRIPTION ??
    "Паломнические поездки в монастыри и святые места России. Сопровождение священника, комфортный транспорт, прозрачная стоимость.",
  region: process.env.NEXT_PUBLIC_SITE_REGION ?? "RU",
  departureCity: process.env.NEXT_PUBLIC_DEPARTURE_CITY ?? "Санкт-Петербург",
  parentOrganization: {
    name: process.env.NEXT_PUBLIC_PARENT_ORG_NAME ?? "Тихвинская епархия",
    url: process.env.NEXT_PUBLIC_PARENT_ORG_URL ?? "https://www.tikhvin-eparhia.ru/",
  },
  contactPhone: process.env.NEXT_PUBLIC_CONTACT_PHONE ?? "+79669334321",
  contactPhoneDisplay: process.env.NEXT_PUBLIC_CONTACT_PHONE_DISPLAY ?? "+7 966 933-43-21",
  contactEmail: process.env.NEXT_PUBLIC_CONTACT_EMAIL ?? "info@tikhvin-palomnik.ru",
  verification: {
    yandex: process.env.NEXT_PUBLIC_YANDEX_VERIFICATION?.trim() || "e79d1ee72d61fee0",
    google:
      process.env.NEXT_PUBLIC_GOOGLE_SITE_VERIFICATION?.trim() ||
      "FvTLt9A-l0U94QIsuHXIBSqygCpr9OAhFsir0BMfbio",
  },
};

export function absoluteUrl(path: string): string {
  const base = siteConfig.url.replace(/\/$/, "");
  const normalized = path.startsWith("/") ? path : `/${path}`;
  return `${base}${normalized}`;
}
