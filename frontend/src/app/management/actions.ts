"use server";

import { revalidatePath } from "next/cache";
import {
  approveManagementReview,
  createManagementCmsBlock,
  createManagementCmsPage,
  bootstrapManagementHomePage,
  createManagementReview,
  createManagementTour,
  deleteManagementCmsBlock,
  deleteManagementCmsPage,
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
  type CmsBlockCreateInput,
  type CmsBlockUpdateInput,
  type CmsPageCreateInput,
  type CmsPageUpdateInput,
  type NewsUpsertInput,
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
};

export async function createReviewAction(input: CreateReviewInput) {
  await createManagementReview({
    tour_id: input.tour_id,
    client_name: input.client_name,
    rating: input.rating,
    text: input.text,
    is_approved: input.is_approved,
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

export async function createCmsPageAction(input: CmsPageCreateInput) {
  const page = await createManagementCmsPage(input);
  revalidateCmsPaths(page.slug);
  revalidatePath(`/management/content/${page.id}`);
  return page;
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

export async function deleteCmsPageAction(formData: FormData) {
  const id = String(formData.get("id") ?? "");
  const slug = String(formData.get("slug") ?? "");
  await deleteManagementCmsPage(id);
  revalidateCmsPaths(slug);
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
