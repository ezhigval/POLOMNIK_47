import { apiGet, apiGetList } from "./client";

export type PublicNewsArticle = {
  id: string;
  slug: string;
  title: string;
  excerpt: string;
  body: string;
  image_url: string;
  published_at: string;
  is_published: boolean;
  sort_order: number;
};

export async function listPublicNews() {
  const body = await apiGetList<PublicNewsArticle>("/news?limit=100");
  return body.data;
}

export async function getPublicNewsBySlug(slug: string) {
  return apiGet<PublicNewsArticle>(`/news/${encodeURIComponent(slug)}`);
}
