"use server";

import { revalidatePath } from "next/cache";
import {
  approveManagementReview,
  createManagementCmsBlock,
  bootstrapManagementHomePage,
  createManagementReview,
  createManagementTour,
  deleteManagementCmsBlock,
  deleteManagementReview,
  deleteManagementTour,
  rejectManagementReview,
  reorderManagementCmsBlocks,
  setManagementReviewReply,
  updateManagementBookingStatus,
  updateManagementCmsBlock,
  updateManagementCmsPage,
  updateManagementTour,
  uploadManagementImage,
  createManagementNews,
  updateManagementNews,
  deleteManagementNews,
  createManagementSMM,
  publishManagementSMM,
  deleteManagementSMM,
  type SMMPostCreateInput,
  updateManagementTelegramSettings,
  updateManagementNotificationSettings,
  updateManagementSiteSettings,
  createManagementRole,
  updateManagementRole,
  deleteManagementRole,
  assignManagementRoleUser,
  sendManagementSupportMessage,
<<<<<<< HEAD
  requestManagementSupportDraft,
=======
  publishManagementLegalDocument,
>>>>>>> f9f53b4 (feat(legal): stages 9-16 reviews, cookie, cabinet, admin, tests)
  type CmsBlockCreateInput,
  type CmsBlockUpdateInput,
  type CmsPageUpdateInput,
  type NewsUpsertInput,
  type SiteSettings,
  type TourUpsertInput,
} from "@/lib/api/management";
import { ApiError } from "@/lib/api/client";

function revalidateCmsPaths(slug?: string) {
  revalidatePath("/");
  revalidatePath("/management/content");
  if (slug && slug !== "home") {
    revalidatePath(`/pages/${slug}`);
  }
  revalidatePath("/pages/[slug]", "page");
}

type CreateReviewInput = {
  tour_id: string;
  client_name: string;
  rating: number;
  text: string;
  is_approved: boolean;
  allow_distribution: boolean;
};

export async function createReviewAction(input: CreateReviewInput) {
  await createManagementReview({
    tour_id: input.tour_id,
    client_name: input.client_name,
    rating: input.rating,
    text: input.text,
    is_approved: input.is_approved,
    allow_distribution: input.allow_distribution,
  });
  revalidateReviewPages();
}

export async function createTourAction(input: TourUpsertInput) {
  await createManagementTour(input);
  revalidatePath("/");
  revalidatePath("/management/tours");
}

export async function uploadTourImageAction(formData: FormData) {
  const file = formData.get("file");
  if (!(file instanceof File) || file.size === 0) {
    throw new Error("Выберите файл изображения");
  }

  try {
    return await uploadManagementImage(file);
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function updateTourAction(id: string, input: TourUpsertInput) {
  await updateManagementTour(id, input);
  revalidatePath("/");
  revalidatePath("/management/tours");
}

export async function deleteTourAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await deleteManagementTour(id);
  revalidatePath("/");
  revalidatePath("/management/tours");
}

export async function updateBookingStatusAction(id: string, status: string) {
  await updateManagementBookingStatus(id, status);
  revalidatePath("/management/bookings");
}

export async function replySupportAction(threadId: string, body: string) {
  await sendManagementSupportMessage(threadId, body);
  revalidatePath("/management/support");
  revalidatePath(`/management/support/${threadId}`);
}

export async function requestSupportDraftAction(threadId: string) {
  try {
    return await requestManagementSupportDraft(threadId);
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function approveReviewAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await approveManagementReview(id);
  revalidateReviewPages();
}

export async function rejectReviewAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await rejectManagementReview(id);
  revalidateReviewPages();
}

export async function deleteReviewAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await deleteManagementReview(id);
  revalidateReviewPages();
}

export async function replyReviewAction(id: string, company_reply: string) {
  await setManagementReviewReply(id, company_reply);
  revalidateReviewPages();
}

function revalidateReviewPages() {
  revalidatePath("/management/reviews");
  revalidatePath("/");
  revalidatePath("/reviews");
  revalidatePath("/tours", "layout");
}

export async function bootstrapHomePageAction() {
  const page = await bootstrapManagementHomePage();
  revalidateCmsPaths(page.slug);
  revalidatePath("/management/content");
  revalidatePath(`/management/content/${page.id}`);
  return page;
}

export async function updateCmsPageAction(id: string, input: CmsPageUpdateInput, slug?: string) {
  const page = await updateManagementCmsPage(id, input);
  revalidateCmsPaths(slug ?? page.slug);
  revalidatePath(`/management/content/${id}`);
  return page;
}

export async function createCmsBlockAction(pageId: string, input: CmsBlockCreateInput, slug?: string) {
  const block = await createManagementCmsBlock(pageId, input);
  revalidateCmsPaths(slug);
  revalidatePath(`/management/content/${pageId}`);
  return block;
}

export async function updateCmsBlockAction(
  id: string,
  input: CmsBlockUpdateInput,
  pageId?: string,
  slug?: string,
) {
  const block = await updateManagementCmsBlock(id, input);
  revalidateCmsPaths(slug);
  if (pageId) {
    revalidatePath(`/management/content/${pageId}`);
  }
  return block;
}

export async function deleteCmsBlockAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  const pageId = String(formData.get("page_id") ?? "");
  const slug = String(formData.get("slug") ?? "");
  await deleteManagementCmsBlock(id);
  revalidateCmsPaths(slug);
  if (pageId) {
    revalidatePath(`/management/content/${pageId}`);
  }
}

export async function reorderCmsBlocksAction(pageId: string, blockIds: string[], slug?: string) {
  const page = await reorderManagementCmsBlocks(pageId, blockIds);
  revalidateCmsPaths(slug ?? page.slug);
  revalidatePath(`/management/content/${pageId}`);
  return page;
}

function revalidateNewsPaths() {
  revalidatePath("/news");
  revalidatePath("/management/news");
}

export async function createNewsAction(input: NewsUpsertInput) {
  await createManagementNews(input);
  revalidateNewsPaths();
}

export async function updateNewsAction(id: string, input: NewsUpsertInput) {
  await updateManagementNews(id, input);
  revalidateNewsPaths();
}

export async function deleteNewsAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await deleteManagementNews(id);
  revalidateNewsPaths();
}

function revalidateSMMPaths() {
  revalidatePath("/management/smm");
  revalidatePath("/news");
}

export async function createSMMPostAction(input: SMMPostCreateInput) {
  await createManagementSMM(input);
  revalidateSMMPaths();
}

export async function publishSMMPostAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await publishManagementSMM(id);
  revalidateSMMPaths();
}

export async function deleteSMMPostAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  await deleteManagementSMM(id);
  revalidateSMMPaths();
}

export async function updateTelegramSettingsAction(input: {
  booking_usernames: string;
  support_usernames: string;
}) {
  try {
    const settings = await updateManagementTelegramSettings(input);
    revalidatePath("/management/settings");
    return settings;
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function updateNotificationSettingsAction(
  events: Record<string, Array<{ channel: string; address: string }>>,
) {
  try {
    const settings = await updateManagementNotificationSettings(events);
    revalidatePath("/management/settings");
    return settings;
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function updateSiteSettingsAction(input: SiteSettings) {
  try {
    const settings = await updateManagementSiteSettings(input);
    revalidatePath("/management/settings");
    return settings;
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function createAdminRoleAction(input: {
  name: string;
  password: string;
  permissions: string[];
}) {
  try {
    const role = await createManagementRole(input);
    revalidatePath("/management/settings");
    return { id: role.id, name: role.name, permissions: role.permissions ?? [] };
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function updateAdminRoleAction(
  id: string,
  input: { password?: string; permissions?: string[] },
) {
  try {
    const role = await updateManagementRole(id, input);
    revalidatePath("/management/settings");
    return role;
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function deleteAdminRoleAction(id: string) {
  try {
    await deleteManagementRole(id);
    revalidatePath("/management/settings");
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function assignAdminRoleUserAction(roleId: string, userId: string) {
  try {
    await assignManagementRoleUser(roleId, userId);
    revalidatePath("/management/settings");
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}

export async function publishLegalDocumentAction(input: {
  type: string;
  version: string;
  title: string;
  content: string;
}) {
  try {
    await publishManagementLegalDocument(input);
    revalidatePath("/management/legal");
    revalidatePath("/legal");
  } catch (error) {
    if (error instanceof ApiError) {
      throw new Error(error.message);
    }
    throw error;
  }
}
