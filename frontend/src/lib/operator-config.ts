/**
 * Единый источник реквизитов оператора ПДн (зеркало backend/internal/legal/operator).
 * После получения настоящих реквизитов заменить placeholders здесь и в backend operator.go.
 */
export const operatorConfig = {
  name: process.env.NEXT_PUBLIC_OPERATOR_NAME ?? "название",
  inn: process.env.NEXT_PUBLIC_OPERATOR_INN ?? "—",
  ogrn: process.env.NEXT_PUBLIC_OPERATOR_OGRN ?? "—",
  legalAddress: process.env.NEXT_PUBLIC_OPERATOR_LEGAL_ADDRESS ?? "—",
  actualAddress: process.env.NEXT_PUBLIC_OPERATOR_ACTUAL_ADDRESS ?? "—",
  email: process.env.NEXT_PUBLIC_OPERATOR_EMAIL ?? process.env.NEXT_PUBLIC_CONTACT_EMAIL ?? "info@tikhvin-palomnik.ru",
  phone: process.env.NEXT_PUBLIC_OPERATOR_PHONE ?? process.env.NEXT_PUBLIC_CONTACT_PHONE_DISPLAY ?? "+7 966 933-43-21",
  website: process.env.NEXT_PUBLIC_SITE_URL ?? "https://tikhvin-palomnik.ru",
  regions: ["Санкт-Петербург", "Ленинградская область", "иные регионы РФ по мере необходимости"],
  publicSiteName: process.env.NEXT_PUBLIC_SITE_NAME ?? "Тихвинский путь",
  publicSiteFull:
    process.env.NEXT_PUBLIC_SITE_FULL_NAME ??
    "РПЦ Тихвинская Епархия. Паломническая служба «Под покровом Божией Матери «Тихвинская»»",
} as const;

/** Типы юридических документов (синхронизированы с backend domain.LegalDocumentType). */
export const legalDocumentTypes = {
  privacyPolicy: "privacy_policy",
  personalData: "personal_data",
  distribution: "distribution",
  marketing: "marketing",
  cookie: "cookie",
  terms: "terms",
  offer: "offer",
} as const;

export type LegalDocumentType = (typeof legalDocumentTypes)[keyof typeof legalDocumentTypes];

/** Маршруты публичных документов на сайте. */
export const legalDocumentPaths: Record<LegalDocumentType, string> = {
  privacy_policy: "/legal/privacy-policy",
  personal_data: "/legal/personal-data-consent",
  distribution: "/legal/distribution-consent",
  marketing: "/legal/marketing-consent",
  cookie: "/legal/cookie-policy",
  terms: "/legal/terms",
  offer: "/legal/offer",
};

export function legalDocumentHref(type: LegalDocumentType): string {
  return legalDocumentPaths[type];
}
