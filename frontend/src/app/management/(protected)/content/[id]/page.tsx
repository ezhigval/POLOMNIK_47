import Link from "next/link";
import { notFound } from "next/navigation";
import { PageEditor } from "@/components/management/cms/page-editor";
import { getManagementCmsPage, listManagementCmsTemplates } from "@/lib/api/management";
import { ApiError } from "@/lib/api/client";

type PageProps = {
  params: Promise<{ id: string }>;
};

export default async function ManagementContentEditorPage({ params }: PageProps) {
  const { id } = await params;

  let page;
  let templates;
  try {
    [page, templates] = await Promise.all([
      getManagementCmsPage(id),
      listManagementCmsTemplates(),
    ]);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      notFound();
    }
    throw error;
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p className="text-sm text-stone-500">
            <Link href="/management" className="hover:text-stone-800">
              ← Обзор
            </Link>
          </p>
          <h1 className="mt-1 text-2xl font-semibold text-stone-900">{page.title}</h1>
        </div>
        {page.is_published ? (
          <Link
            href={page.slug === "home" ? "/" : `/pages/${page.slug}`}
            className="btn-secondary text-sm"
            target="_blank"
          >
            Открыть на сайте
          </Link>
        ) : null}
      </div>

      <PageEditor page={page} templates={templates} />
    </div>
  );
}
