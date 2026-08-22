import { CreateNewsForm } from "@/components/management/create-news-form";
import { EditNewsForm } from "@/components/management/edit-news-form";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { deleteNewsAction } from "@/app/management/actions";
import { listManagementNews } from "@/lib/api/management";
import { formatNewsDate } from "@/lib/news";

export default async function ManagementNewsPage() {
  const articles = await listManagementNews();

  return (
    <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
      <ManagementPanel title="Новости" description={`${articles.length} статей`}>
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Статья</ManagementTh>
            <ManagementTh>Дата</ManagementTh>
            <ManagementTh>Статус</ManagementTh>
            <ManagementTh />
          </ManagementTableHead>
          <tbody>
            {articles.length === 0 ? (
              <ManagementEmptyRow colSpan={4}>Статей пока нет.</ManagementEmptyRow>
            ) : (
              articles.map((article) => (
                <tr key={article.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4">
                    <div className="font-medium text-stone-900">{article.title}</div>
                    <div className="text-stone-500">{article.slug}</div>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap">{formatNewsDate(article.published_at)}</td>
                  <td className="px-4 py-4">
                    <StatusBadge variant={article.is_published ? "success" : "neutral"}>
                      {article.is_published ? "Опубликована" : "Черновик"}
                    </StatusBadge>
                  </td>
                  <td className="px-4 py-4">
                    <div className="space-y-2">
                      <EditNewsForm article={article} />
                      <form action={deleteNewsAction}>
                        <input type="hidden" name="id" value={article.id} />
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

      <CreateNewsForm />
    </div>
  );
}
