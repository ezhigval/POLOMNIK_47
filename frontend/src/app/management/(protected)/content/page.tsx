import Link from "next/link";
import { deleteCmsPageAction } from "@/app/management/actions";
import { BootstrapHomeButton } from "@/components/management/cms/bootstrap-home-button";
import { CreatePageForm } from "@/components/management/cms/create-page-form";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { listManagementCmsPagesOrEmpty } from "@/lib/api/management";

export default async function ManagementContentPage() {
  const { pages, unavailable } = await listManagementCmsPagesOrEmpty();
  const hasHome = pages.some((page) => page.slug === "home");

  return (
    <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
      <div className="space-y-6">
        {!unavailable && !hasHome ? <BootstrapHomeButton /> : null}

        <ManagementPanel title="Страницы" description={unavailable ? "недоступно" : `${pages.length} в CMS`}>
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Название</ManagementTh>
            <ManagementTh>Slug</ManagementTh>
            <ManagementTh>Path</ManagementTh>
            <ManagementTh>Статус</ManagementTh>
            <ManagementTh />
          </ManagementTableHead>
          <tbody>
            {pages.length === 0 ? (
              <ManagementEmptyRow colSpan={5}>
                {unavailable ? "Список страниц сейчас недоступен." : "Страниц пока нет."}
              </ManagementEmptyRow>
            ) : (
              pages.map((page) => (
                <tr key={page.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4">
                    <div className="font-medium text-stone-900">{page.title}</div>
                  </td>
                  <td className="px-4 py-4 text-stone-600">{page.slug}</td>
                  <td className="px-4 py-4 text-stone-600">{page.path}</td>
                  <td className="px-4 py-4">
                    <StatusBadge variant={page.is_published ? "success" : "neutral"}>
                      {page.is_published ? "Опубликована" : "Черновик"}
                    </StatusBadge>
                  </td>
                  <td className="px-4 py-4">
                    <div className="space-y-2">
                      <Link
                        href={`/management/content/${page.id}`}
                        className="text-sm font-medium text-brand-800 hover:text-brand-900"
                      >
                        Редактировать
                      </Link>
                      <form action={deleteCmsPageAction}>
                        <input type="hidden" name="id" value={page.id} />
                        <input type="hidden" name="slug" value={page.slug} />
                        <button type="submit" className="btn-danger">
                          Удалить
                        </button>
                      </form>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>
      </div>

      <CreatePageForm />
    </div>
  );
}
