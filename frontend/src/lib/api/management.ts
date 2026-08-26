import "server-only";

import { cookies } from "next/headers";
import { apiUrl, requestJson, ApiError, type DataEnvelope, type ListEnvelope } from "./client";
import {
  ADMIN_SESSION_COOKIE,
  isManagementJwt,
  verifyAdminSessionValue,
} from "@/lib/auth/admin-session";
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
  payment_status: string;
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
  is_regular: boolean;
  overbooking_enabled: boolean;
};

export function isManagementConfigured(): boolean {
  return Boolean(process.env.ADMIN_TOKEN);
}

async function managementAuthHeaders(): Promise<Record<string, string>> {
  const cookieStore = await cookies();
  const session = cookieStore.get(ADMIN_SESSION_COOKIE)?.value;
  if (!session) {
    throw new ApiError(401, "UNAUTHORIZED", "Нужна авторизация администратора");
  }
  if (isManagementJwt(session)) {
    return { "X-Admin-Session": session };
  }
  const adminToken = process.env.ADMIN_TOKEN;
  if (adminToken && verifyAdminSessionValue(session, adminToken)) {
    return { "X-Admin-Token": adminToken };
  }
  throw new ApiError(401, "UNAUTHORIZED", "Нужна авторизация администратора");
}

async function managementRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const auth = await managementAuthHeaders();
  return requestJson<T>(apiUrl(`/management${path}`), {
    cache: "no-store",
    ...init,
    headers: {
      ...auth,
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
  const auth = await managementAuthHeaders();

  const body = new FormData();
  body.append("file", file);

  const payload = await requestJson<DataEnvelope<{ url: string; path: string }>>(apiUrl("/management/uploads"), {
    method: "POST",
    cache: "no-store",
    headers: {
      ...auth,
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
  messenger_adapter?: string;
  publisher_adapter?: string;
  ai_adapter?: string;
  payment_adapter?: string;
  telegram_configured: boolean;
  messenger_configured?: boolean;
  publisher_configured?: boolean;
  ai_configured?: boolean;
  payment_configured?: boolean;
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
  latency?: {
    last_ms: number;
    avg_ms: number;
    requests: number;
  };
  last_backup?: {
    at?: string;
    file?: string;
    bytes?: number;
    offsite: boolean;
    offsite_error?: string;
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

export async function listManagementBookings(params?: {
  status?: string;
  date_from?: string;
  date_to?: string;
  tour_id?: string;
  page?: number;
  limit?: number;
}) {
  const search = new URLSearchParams();
  if (params?.status) search.set("status", params.status);
  if (params?.date_from) search.set("date_from", params.date_from);
  if (params?.date_to) search.set("date_to", params.date_to);
  if (params?.tour_id) search.set("tour_id", params.tour_id);
  if (params?.page) search.set("page", String(params.page));
  if (params?.limit) search.set("limit", String(params.limit));
  const query = search.toString();
  const body = await managementRequest<ListEnvelope<ManagementBooking>>(
    `/bookings${query ? `?${query}` : ""}`,
  );
  return body.data;
}

export async function exportManagementBookingsCSV(params: {
  status?: string;
  date_from?: string;
  date_to?: string;
  tour_id?: string;
}) {
  const auth = await managementAuthHeaders();
  const search = new URLSearchParams();
  search.set("format", "csv");
  if (params.status) search.set("status", params.status);
  if (params.date_from) search.set("date_from", params.date_from);
  if (params.date_to) search.set("date_to", params.date_to);
  if (params.tour_id) search.set("tour_id", params.tour_id);
  const response = await fetch(apiUrl(`/management/bookings?${search.toString()}`), {
    headers: auth,
    cache: "no-store",
  });
  return response;
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

export async function updateManagementBookingPaymentStatus(id: string, paymentStatus: string) {
  const body = await managementRequest<DataEnvelope<ManagementBooking>>(
    `/bookings/${id}/payment-status`,
    {
      method: "PATCH",
      body: JSON.stringify({ payment_status: paymentStatus }),
    },
  );
  return body.data;
}

export type ManagementSupportMessage = {
  id: string;
  sender_type: string;
  body: string;
  created_at: string;
};

export type ManagementSupportThread = {
  id: string;
  user_id: string;
  subject: string;
  status: string;
  messages?: ManagementSupportMessage[];
  created_at?: string;
  updated_at: string;
};

export async function listManagementSupportThreads() {
  const body = await managementRequest<ListEnvelope<ManagementSupportThread>>("/support");
  return body.data;
}

export async function getManagementSupportThread(id: string) {
  const body = await managementRequest<DataEnvelope<ManagementSupportThread>>(`/support/${id}`);
  return body.data;
}

export type SupportDraft = {
  configured: boolean;
  escalate: boolean;
  draft: string;
  note: string;
};

export async function requestManagementSupportDraft(id: string) {
  const body = await managementRequest<DataEnvelope<SupportDraft>>(`/support/${id}/draft`, {
    method: "POST",
  });
  return body.data;
}

export type MetricsDigest = {
  configured: boolean;
  bookings_by_status: Record<string, number>;
  active_tours: number;
  open_support_threads: number;
  outbox_pending: number;
  outbox_failed: number;
  summary?: string;
};

export async function getManagementMetricsDigest() {
  const body = await managementRequest<DataEnvelope<MetricsDigest>>("/ai/metrics-digest");
  return body.data;
}

export type WatchdogReport = {
  configured: boolean;
  at: string;
  database: string;
  disk_path?: string;
  disk_used_bytes: number;
  disk_total_bytes: number;
  disk_percent: number;
  outbox_pending: number;
  outbox_failed: number;
  status_5xx: number;
  backup_at?: string;
  backup_overdue: boolean;
  restart_attempted: boolean;
  summary?: string;
};

export async function getManagementWatchdog() {
  const body = await managementRequest<DataEnvelope<WatchdogReport>>("/watchdog");
  return body.data;
}

export async function sendManagementSupportMessage(id: string, bodyText: string) {
  const body = await managementRequest<DataEnvelope<ManagementSupportThread>>(
    `/support/${id}/messages`,
    {
      method: "POST",
      body: JSON.stringify({ body: bodyText }),
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
  allow_distribution: boolean;
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

export type CmsPageUpdateInput = {
  title?: string;
  path?: string;
  meta_title?: string;
  meta_description?: string;
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
  is_pinned: boolean;
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
  is_pinned: boolean;
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

export async function setManagementNewsPinned(id: string, isPinned: boolean) {
  const body = await managementRequest<DataEnvelope<ManagementNewsArticle>>(`/news/${id}/pin`, {
    method: "PATCH",
    body: JSON.stringify({ is_pinned: isPinned }),
  });
  return body.data;
}

export type ManagementSMMResult = {
  channel: string;
  ok: boolean;
  error?: string;
  attempted_at?: string;
};

export type ManagementSMMPost = {
  id: string;
  title: string;
  body: string;
  url: string;
  publish_at: string;
  channels: string[];
  published_at?: string;
  results: ManagementSMMResult[];
  created_at: string;
};

export type SMMPostCreateInput = {
  title: string;
  body: string;
  url: string;
  publish_at: string;
  channels: string[];
};

export async function listManagementSMM() {
  const body = await managementRequest<ListEnvelope<ManagementSMMPost>>("/smm?limit=100");
  return body.data;
}

export async function createManagementSMM(input: SMMPostCreateInput) {
  const body = await managementRequest<DataEnvelope<ManagementSMMPost>>("/smm", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function publishManagementSMM(id: string) {
  const body = await managementRequest<DataEnvelope<ManagementSMMPost>>(`/smm/${id}/publish`, {
    method: "POST",
  });
  return body.data;
}

export async function deleteManagementSMM(id: string) {
  await managementRequest<void>(`/smm/${id}`, { method: "DELETE" });
}

export type NotificationSettings = {
  channels: Array<{ id: string; configured: boolean; label: string }>;
  events: Array<{
    kind: string;
    title: string;
    recipients: Array<{
      channel: string;
      address: string;
      event: string;
      ready: boolean;
      status: string;
      label: string;
    }>;
  }>;
};

export type SiteSettings = {
  site_name: string;
  full_name: string;
  tagline: string;
  description: string;
  region: string;
  departure_city: string;
  parent_org_name: string;
  parent_org_url: string;
  contact_phone: string;
  contact_phone_display: string;
  contact_email: string;
  mail_forward_to: string;
};

export type AdminRole = {
  id: string;
  name: string;
  permissions: string[];
  created_at?: string;
  updated_at?: string;
};

export type AdminRoleTemplate = {
  id: string;
  label: string;
  role_name: string;
  permissions: string[];
};

export type ManagementSession = {
  full_admin: boolean;
  role_id?: string;
  role_name?: string;
  permissions: string[];
};

export async function getManagementSession() {
  const body = await managementRequest<DataEnvelope<ManagementSession>>("/session");
  return body.data;
}

export async function getManagementNotificationSettings() {
  const body = await managementRequest<DataEnvelope<NotificationSettings>>("/notification-settings");
  return body.data;
}

export async function updateManagementNotificationSettings(
  events: Record<string, Array<{ channel: string; address: string }>>,
) {
  const body = await managementRequest<DataEnvelope<NotificationSettings>>("/notification-settings", {
    method: "PATCH",
    body: JSON.stringify({ events }),
  });
  return body.data;
}

export async function getManagementSiteSettings() {
  const body = await managementRequest<DataEnvelope<SiteSettings>>("/site-settings");
  return body.data;
}

export async function updateManagementSiteSettings(input: SiteSettings) {
  const body = await managementRequest<DataEnvelope<SiteSettings>>("/site-settings", {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function listManagementRoles() {
  const body = await managementRequest<DataEnvelope<AdminRole[]>>("/roles");
  return body.data ?? [];
}

export async function listManagementRoleTemplates() {
  const body = await managementRequest<DataEnvelope<AdminRoleTemplate[]>>("/roles/templates");
  return body.data ?? [];
}

export async function createManagementRole(input: {
  name: string;
  password: string;
  permissions: string[];
}) {
  const body = await managementRequest<DataEnvelope<AdminRole>>("/roles", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function updateManagementRole(
  id: string,
  input: { password?: string; permissions?: string[] },
) {
  const body = await managementRequest<DataEnvelope<AdminRole>>(`/roles/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function deleteManagementRole(id: string) {
  await managementRequest<void>(`/roles/${id}`, { method: "DELETE" });
}

export async function assignManagementRoleUser(roleId: string, userId: string) {
  await managementRequest<DataEnvelope<{ status: string }>>(`/roles/${roleId}/assignments`, {
    method: "POST",
    body: JSON.stringify({ user_id: userId }),
  });
}

export type ManagementLegalDocument = {
  id: string;
  type: string;
  version: string;
  title: string;
  published_at: string;
  updated_at: string;
  is_active: boolean;
  content?: string;
};

export type ManagementConsent = {
  id: string;
  user_id?: string | null;
  request_id?: string | null;
  consent_type: string;
  document_id: string;
  document_version: string;
  accepted_at: string;
};

export async function listManagementLegalDocuments(type?: string) {
  const query = type ? `?type=${encodeURIComponent(type)}` : "";
  const body = await managementRequest<DataEnvelope<ManagementLegalDocument[]>>(
    `/legal/documents${query}`,
  );
  return body.data;
}

export async function publishManagementLegalDocument(input: {
  type: string;
  version: string;
  title: string;
  content: string;
}) {
  const body = await managementRequest<DataEnvelope<ManagementLegalDocument>>("/legal/documents", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return body.data;
}

export async function listManagementConsents(page = 1, limit = 50) {
  return managementRequest<ListEnvelope<ManagementConsent>>(
    `/consents?page=${page}&limit=${limit}`,
  );
}
