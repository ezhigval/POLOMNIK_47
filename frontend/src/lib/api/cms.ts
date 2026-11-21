import { ApiError, apiGet, apiGetList } from "./client";

export type CmsBlockType =
  | "hero"
  | "about"
  | "why_us"
  | "how_it_works"
  | "faq"
  | "cta"
  | "rich_text"
  | "popular_destinations"
  | "testimonials"
  | "featured_route"
  | "information_hero"
  | "holy_places"
  | "service_profile"
  | "service_team"
  | "legal_donation"
  | "labor_trips";

export type HeroBlockContent = {
  eyebrow?: string;
  title?: string;
  subtitle?: string;
  primaryCta?: string;
  primaryHref?: string;
  secondaryCta?: string;
  secondaryHref?: string;
  stats?: { value: string; label: string }[];
};

export type AboutBlockContent = {
  eyebrow?: string;
  title?: string;
  paragraphs?: string[];
  highlights?: string[];
  showContacts?: boolean;
  stats?: { value: string; label: string }[];
};

export type WhyUsBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  items?: { title: string; description: string; icon: string }[];
  stats?: { value: string; label: string }[];
};

export type HowItWorksBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  steps?: { title: string; description: string }[];
  ctaLabel?: string;
  ctaHref?: string;
};

export type FaqBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  items?: { question: string; answer: string }[];
};

export type CtaBlockContent = {
  title?: string;
  subtitle?: string;
  button?: string;
  href?: string;
};

export type RichTextBlockContent = {
  eyebrow?: string;
  title?: string;
  body?: string;
};

export type FeaturedRouteBlockContent = {
  eyebrow?: string;
  title?: string;
  parentRoute?: string;
  duration?: string;
  region?: string;
  lead?: string;
  body?: string;
  ctaLabel?: string;
  ctaHref?: string;
  secondaryCta?: string;
  secondaryHref?: string;
  days?: { title: string; points: string[] }[];
};

export type InformationHeroBlockContent = Pick<HeroBlockContent, "eyebrow" | "title" | "subtitle"> & {
  primaryCta?: string;
  primaryHref?: string;
  secondaryCta?: string;
  secondaryHref?: string;
};

export type HolyPlacesBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  items?: {
    name: string;
    region: string;
    description: string;
    address: string;
    directions?: string;
    latitude: string;
    longitude: string;
    mapUrl: string;
    phone: string;
    email: string;
    phoneLabel?: string;
    websiteUrl: string;
    bookingHref: string;
    bookingLabel: string;
  }[];
};

export type ServiceProfileBlockContent = {
  eyebrow?: string;
  title?: string;
  paragraphs?: string[];
  links?: { label: string; href: string }[];
};

export type ServiceTeamBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  members?: { photo?: string; name: string; role: string; bio: string; phone: string; email: string }[];
};

export type LegalDonationBlockContent = {
  eyebrow?: string;
  title?: string;
  legalName?: string;
  legalAddress?: string;
  inn?: string;
  kpp?: string;
  ogrn?: string;
  donationRecipient?: string;
  donationBankName?: string;
  donationAccount?: string;
  donationCorrespondentAccount?: string;
  donationBik?: string;
  donationPurpose?: string;
  donationUrl?: string;
};

export type LaborTripsBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  items?: { name: string; place: string; dates: string; description: string; bookingHref: string }[];
};

export type CmsBlockContent =
  | HeroBlockContent
  | AboutBlockContent
  | WhyUsBlockContent
  | HowItWorksBlockContent
  | FaqBlockContent
  | CtaBlockContent
  | RichTextBlockContent
  | FeaturedRouteBlockContent
  | InformationHeroBlockContent
  | HolyPlacesBlockContent
  | ServiceProfileBlockContent
  | ServiceTeamBlockContent
  | LegalDonationBlockContent
  | LaborTripsBlockContent
  | Record<string, never>;

export type CmsBlock = {
  id: string;
  page_id?: string;
  type: CmsBlockType;
  sort_order: number;
  content: Record<string, unknown>;
  is_visible: boolean;
  created_at?: string;
  updated_at?: string;
};

export type CmsPage = {
  id: string;
  slug: string;
  title: string;
  path: string;
  meta_title?: string;
  meta_description?: string;
  is_published: boolean;
  created_at?: string;
  updated_at?: string;
  blocks: CmsBlock[];
};

export type CmsBlockTemplate = {
  type: CmsBlockType;
  label: string;
  content: Record<string, unknown>;
};

export async function getPublishedPage(slug: string): Promise<CmsPage | null> {
  try {
    return await apiGet<CmsPage>(`/pages/${encodeURIComponent(slug)}`);
  } catch (error) {
    if (error instanceof ApiError && (error.status === 404 || error.status >= 500)) {
      return null;
    }
    throw error;
  }
}

export async function listPublishedPages(): Promise<CmsPage[]> {
  try {
    const body = await apiGetList<CmsPage>("/pages");
    return body.data ?? [];
  } catch (error) {
    if (error instanceof ApiError && (error.status === 404 || error.status >= 500)) {
      return [];
    }
    throw error;
  }
}

export function cmsPublicPath(page: Pick<CmsPage, "slug" | "path">): string | null {
  if (page.slug === "home" || page.path === "/") {
    return null;
  }
  if (page.path?.startsWith("/") && page.path !== "/") {
    return page.path;
  }
  return `/pages/${page.slug}`;
}
