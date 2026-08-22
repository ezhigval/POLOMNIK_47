import "server-only";

import { apiUrl, requestJson, ApiError, type DataEnvelope, type ListEnvelope } from "./client";
import type { CmsBlock, CmsBlockTemplate, CmsBlockType, CmsPage } from "./cms";
import type { Tour } from "./tours";

export type ManagementTour = Tour & {
  is_active: boolean;
  overbooking_enabled: boolean;
};

export type ManagementBooking = {
  id: string;
  tour_id: string;
  name: string;
  phone: string;
  email: string;
  people_count: number;
  status: string;
  total_price: number;
  comment: string;
  overbooked: boolean;
  source: string;
  created_at: string;
  updated_at: string;
};

export type ManagementReview = {
  id: string;
  tour_id: string;
  client_name: string;
  rating: number;
  text: string;
  company_reply: string;
  company_replied_at?: string | null;
  is_approved: boolean;
  created_at: string;
};

export type TourUpsertInput = {
  slug: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  date_start: string;
  date_end: string;
  slots_total: number;
  slots_left: number;
  location: string;
  images: string[];
  is_active: boolean;
  is_hot: boolean;
  overbooking_enabled: boolean;
};

export function isManagementConfigured(): boolean {
  return Boolean(process.env.ADMIN_TOKEN);
}

async function managementRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const adminToken = process.env.ADMIN_TOKEN;
  if (!adminToken) {
    throw new ApiError(503, "SERVICE_UNAVAILABLE", "Management API is not configured");
  }

  return requestJson<T>(apiUrl(`/management${path}`), {
    cache: "no-store",
    ...init,
    headers: {
      "X-Admin-Token": adminToken,
      ...(init?.headers ?? {}),
    },
  });
}

export async function listManagementTours() {
  const body = await managementRequest<ListEnvelope<ManagementTour>>("/tours");
  return body.data;
}

export async function createManagementTour(input: TourUpsertInput) {
  const body = await managementRequest<DataEnvelope<ManagementTour>>("/tours", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function uploadManagementImage(file: File) {
  const adminToken = process.env.ADMIN_TOKEN;
  if (!adminToken) {
    throw new ApiError(503, "SERVICE_UNAVAILABLE", "Management API is not configured");
  }

  const body = new FormData();
  body.append("file", file);

  const payload = await requestJson<DataEnvelope<{ url: string; path: string }>>(apiUrl("/management/uploads"), {
    method: "POST",
    cache: "no-store",
    headers: {
      "X-Admin-Token": adminToken,
    },
    body,
  });
  return payload.data;
}

export async function updateManagementTour(id: string, input: TourUpsertInput) {
  const body = await managementRequest<DataEnvelope<ManagementTour>>(`/tours/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}

export type ManagementSystemInfo = {
  crm_adapter: string;
  accounting_adapter: string;
  notification_adapter: string;
  telegram_configured: boolean;
  bitrix_configured: boolean;
  onec_configured: boolean;
  outbox: {
    pending: number;
    failed: number;
    processed: number;
    oldest_pending_at?: string;
    latest_failed_at?: string;
    latest_failed_error?: string;
  };
};

export async function getManagementSystemInfo() {
  const body = await managementRequest<DataEnvelope<ManagementSystemInfo>>("/system-info");
  return body.data;
}

export async function listManagementIntegrationReferences(params?: {
  external_system?: string;
  local_entity_type?: string;
  sync_status?: string;
}) {
  const search = new URLSearchParams();
  if (params?.external_system) search.set("external_system", params.external_system);
  if (params?.local_entity_type) search.set("local_entity_type", params.local_entity_type);
  if (params?.sync_status) search.set("sync_status", params.sync_status);

  const query = search.toString();
  const body = await managementRequest<ListEnvelope<ManagementIntegrationReference>>(
    `/integration-references${query ? `?${query}` : ""}`,
  );
  return body.data;
}

export async function listManagementOutboxEvents(params?: {
  status?: string;
  entity_type?: string;
  event_type?: string;
}) {
  const search = new URLSearchParams();
  if (params?.status) search.set("status", params.status);
  if (params?.entity_type) search.set("entity_type", params.entity_type);
  if (params?.event_type) search.set("event_type", params.event_type);

  const query = search.toString();
  const body = await managementRequest<ListEnvelope<ManagementOutboxEvent>>(
    `/outbox-events${query ? `?${query}` : ""}`,
  );
  return body.data;
}

export type ManagementOutboxEvent = {
  id: string;
  event_type: string;
  entity_type: string;
  entity_id: string;
  payload: Record<string, unknown>;
  status: string;
  attempts: number;
  last_error: string;
  created_at: string;
  updated_at: string;
};

export type ManagementIntegrationReference = {
  id: string;
  local_entity_type: string;
  local_entity_id: string;
  external_system: string;
  external_entity_type: string;
  external_entity_id: string;
  sync_status: string;
  last_sync_at: string | null;
  last_error: string;
  created_at: string;
  updated_at: string;
};

export async function deleteManagementTour(id: string) {
  await managementRequest<void>(`/tours/${id}`, { method: "DELETE" });
}

export async function listManagementBookings() {
  const body = await managementRequest<ListEnvelope<ManagementBooking>>("/bookings");
  return body.data;
}

export async function updateManagementBookingStatus(id: string, status: string) {
  const body = await managementRequest<DataEnvelope<ManagementBooking>>(
    `/bookings/${id}/status`,
    {
      method: "PATCH",
      body: JSON.stringify({ status }),
    },
  );
  return body.data;
}

export type CreateReviewInput = {
  tour_id: string;
  client_name: string;
  rating: number;
  text: string;
  is_approved: boolean;
};

export async function createManagementReview(input: CreateReviewInput) {
  const body = await managementRequest<DataEnvelope<ManagementReview>>("/reviews", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function listManagementReviews() {
  const body = await managementRequest<ListEnvelope<ManagementReview>>("/reviews");
  return body.data;
}

export async function approveManagementReview(id: string) {
  const body = await managementRequest<DataEnvelope<ManagementReview>>(
    `/reviews/${id}/approve`,
    { method: "PATCH" },
  );
  return body.data;
}

export async function rejectManagementReview(id: string) {
  const body = await managementRequest<DataEnvelope<ManagementReview>>(
    `/reviews/${id}/reject`,
    { method: "PATCH" },
  );
  return body.data;
}

export async function setManagementReviewReply(id: string, company_reply: string) {
  const body = await managementRequest<DataEnvelope<ManagementReview>>(
    `/reviews/${id}/reply`,
    {
      method: "PATCH",
      body: JSON.stringify({ company_reply }),
    },
  );
  return body.data;
}

export async function deleteManagementReview(id: string) {
  await managementRequest<void>(`/reviews/${id}`, { method: "DELETE" });
}

export type CmsPageCreateInput = {
  slug: string;
  title: string;
  path: string;
  is_published: boolean;
};

export type CmsPageUpdateInput = {
  title?: string;
  path?: string;
  is_published?: boolean;
};

export type CmsBlockCreateInput = {
  type: CmsBlockType;
  content?: Record<string, unknown>;
  is_visible?: boolean;
};

export type CmsBlockUpdateInput = {
  content?: Record<string, unknown>;
  is_visible?: boolean;
  sort_order?: number;
};

export async function listManagementCmsPages() {
  const body = await managementRequest<ListEnvelope<CmsPage>>("/cms/pages");
  return body.data ?? [];
}

export async function listManagementCmsPagesOrEmpty(): Promise<{ pages: CmsPage[]; unavailable: boolean }> {
  try {
    return { pages: await listManagementCmsPages(), unavailable: false };
  } catch (error) {
    if (error instanceof ApiError) {
      return { pages: [], unavailable: true };
    }
    throw error;
  }
}

export async function getManagementCmsPage(id: string) {
  const body = await managementRequest<DataEnvelope<CmsPage>>(`/cms/pages/${id}`);
  return body.data;
}

export async function createManagementCmsPage(input: CmsPageCreateInput) {
  const body = await managementRequest<DataEnvelope<CmsPage>>("/cms/pages", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function bootstrapManagementHomePage() {
  const body = await managementRequest<DataEnvelope<CmsPage>>("/cms/pages/bootstrap-home", {
    method: "POST",
  });
  return body.data;
}

export async function updateManagementCmsPage(id: string, input: CmsPageUpdateInput) {
  const body = await managementRequest<DataEnvelope<CmsPage>>(`/cms/pages/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function deleteManagementCmsPage(id: string) {
  await managementRequest<void>(`/cms/pages/${id}`, { method: "DELETE" });
}

export async function listManagementCmsTemplates() {
  const body = await managementRequest<DataEnvelope<CmsBlockTemplate[]>>("/cms/templates");
  return body.data;
}

export async function createManagementCmsBlock(pageId: string, input: CmsBlockCreateInput) {
  const body = await managementRequest<DataEnvelope<CmsBlock>>(`/cms/pages/${pageId}/blocks`, {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function updateManagementCmsBlock(id: string, input: CmsBlockUpdateInput) {
  const body = await managementRequest<DataEnvelope<CmsBlock>>(`/cms/blocks/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function deleteManagementCmsBlock(id: string) {
  await managementRequest<void>(`/cms/blocks/${id}`, { method: "DELETE" });
}

export async function reorderManagementCmsBlocks(pageId: string, blockIds: string[]) {
  const body = await managementRequest<DataEnvelope<CmsPage>>(`/cms/pages/${pageId}/blocks/reorder`, {
    method: "POST",
    body: JSON.stringify({ block_ids: blockIds }),
  });
  return body.data;
}

export type ManagementNewsArticle = {
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

export type NewsUpsertInput = {
  slug: string;
  title: string;
  excerpt: string;
  body: string;
  image_url: string;
  published_at: string;
  is_published: boolean;
  sort_order: number;
};

export async function listManagementNews() {
  const body = await managementRequest<ListEnvelope<ManagementNewsArticle>>("/news?limit=100");
  return body.data;
}

export async function createManagementNews(input: NewsUpsertInput) {
  const body = await managementRequest<DataEnvelope<ManagementNewsArticle>>("/news", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function updateManagementNews(id: string, input: NewsUpsertInput) {
  const body = await managementRequest<DataEnvelope<ManagementNewsArticle>>(`/news/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function deleteManagementNews(id: string) {
  await managementRequest<void>(`/news/${id}`, { method: "DELETE" });
}

export type ManagementTelegramRecipient = {
  username: string;
  kind: "booking" | "support";
  chat_bound: boolean;
  status: string;
};

export type ManagementTelegramSettings = {
  booking_usernames: string;
  support_usernames: string;
  recipients: ManagementTelegramRecipient[];
};

export async function getManagementTelegramSettings() {
  const body = await managementRequest<DataEnvelope<ManagementTelegramSettings>>("/telegram-settings");
  return body.data;
}

function normalizeTelegramSettings(
  settings: ManagementTelegramSettings | null | undefined,
): ManagementTelegramSettings | null {
  if (!settings) {
    return null;
  }
  return {
    booking_usernames: settings.booking_usernames ?? "",
    support_usernames: settings.support_usernames ?? "",
    recipients: Array.isArray(settings.recipients) ? settings.recipients : [],
  };
}

export async function getManagementTelegramSettingsOrEmpty(): Promise<{
  settings: ManagementTelegramSettings | null;
  unavailable: boolean;
}> {
  try {
    return {
      settings: normalizeTelegramSettings(await getManagementTelegramSettings()),
      unavailable: false,
    };
  } catch {
    return { settings: null, unavailable: true };
  }
}

export async function updateManagementTelegramSettings(input: {
  booking_usernames: string;
  support_usernames: string;
}) {
  const body = await managementRequest<DataEnvelope<ManagementTelegramSettings>>("/telegram-settings", {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}
