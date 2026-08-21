import { ApiError, apiGet } from "./client";

export type CmsBlockType =
  | "hero"
  | "about"
  | "why_us"
  | "how_it_works"
  | "faq"
  | "cta"
  | "rich_text"
  | "popular_destinations"
  | "testimonials";

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
};

export type WhyUsBlockContent = {
  eyebrow?: string;
  title?: string;
  description?: string;
  items?: { title: string; description: string; icon: string }[];
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

export type CmsBlockContent =
  | HeroBlockContent
  | AboutBlockContent
  | WhyUsBlockContent
  | HowItWorksBlockContent
  | FaqBlockContent
  | CtaBlockContent
  | RichTextBlockContent
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
    return await apiGet<CmsPage>(`/api/v1/pages/${encodeURIComponent(slug)}`);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      return null;
    }
    throw error;
  }
}
