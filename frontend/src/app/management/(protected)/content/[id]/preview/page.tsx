import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { CmsPageRenderer } from "@/components/cms/cms-page-renderer";
import { InformationSpace, informationBlockTypes } from "@/components/information-space";
import { getManagementCmsPage } from "@/lib/api/management";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";
import { ApiError } from "@/lib/api/client";

export const dynamic = "force-dynamic";
export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default async function ManagementCmsLivePreview({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ space?: string }>;
}) {
  if (!(await canAccessManagementPage([PERM.content]))) {
    notFound();
  }

  const [{ id }, query] = await Promise.all([params, searchParams]);
  let page;
  try {
    page = await getManagementCmsPage(id);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) notFound();
    throw error;
  }

  const isHome = page.slug === "home";
  const guide = query.space !== "trips";
  const blocks = page.blocks ?? [];

  return (
    <main className="min-w-0 bg-stone-50" data-preview="cms" data-home-mode={isHome ? (guide ? "guide" : "trips") : undefined}>
      <div className="border-b border-amber-200 bg-amber-50 px-4 py-2 text-center text-xs font-medium text-amber-900">
        Предпросмотр сохранённых данных: изменения видны только после сохранения. Публикация не требуется.
      </div>
      {isHome && guide ? (
        <InformationSpace blocks={blocks} fallback={false} />
      ) : (
        <CmsPageRenderer blocks={isHome ? blocks.filter((block) => !informationBlockTypes.has(block.type)) : blocks} />
      )}
    </main>
  );
}
