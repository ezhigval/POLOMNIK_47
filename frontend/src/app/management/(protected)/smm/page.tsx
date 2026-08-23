import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { CreateSMMForm } from "@/components/management/create-smm-form";
import { deleteSMMPostAction, publishSMMPostAction } from "@/app/management/actions";
import { listManagementSMM } from "@/lib/api/management";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

export default async function ManagementSMMPage() {
  if (!(await canAccessManagementPage([PERM.content]))) {
    return <ManagementNoAccess />;
  }
  const posts = await listManagementSMM();

  return (
    <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
      <ManagementPanel
        title="Контент-план"
        description="Материал, слот и список publisher. Падение одного канала не откатывает остальные. На проде PublisherPort по умолчанию noop."
      >
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Материал</ManagementTh>
            <ManagementTh>Слот</ManagementTh>
            <ManagementTh>Каналы</ManagementTh>
            <ManagementTh />
          </ManagementTableHead>
          <tbody>
            {posts.length === 0 ? (
              <ManagementEmptyRow colSpan={4}>Материалов пока нет.</ManagementEmptyRow>
            ) : (
              posts.map((post) => (
                <tr key={post.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4">
                    <div className="font-medium text-stone-900">{post.title}</div>
                    <div className="line-clamp-2 text-stone-500">{post.body}</div>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm">{post.publish_at}</td>
                  <td className="px-4 py-4 text-sm">
                    <div>{post.channels.join(", ")}</div>
                    {post.published_at ? (
                      <div className="mt-2 space-y-1">
                        {post.results.map((item) => (
                          <StatusBadge key={item.channel} variant={item.ok ? "success" : "warning"}>
                            {item.channel}: {item.ok ? "ok" : item.error || "ошибка"}
                          </StatusBadge>
                        ))}
                      </div>
                    ) : (
                      <StatusBadge variant="neutral">в плане</StatusBadge>
                    )}
                  </td>
                  <td className="px-4 py-4 space-y-2">
                    {!post.published_at ? (
                      <form action={publishSMMPostAction}>
                        <input type="hidden" name="id" value={post.id} />
                        <button type="submit" className="btn-primary">
                          Опубликовать
                        </button>
                      </form>
                    ) : null}
                    <form action={deleteSMMPostAction}>
                      <input type="hidden" name="id" value={post.id} />
                      <button type="submit" className="btn-danger">
                        Удалить
                      </button>
                    </form>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>
      <CreateSMMForm />
    </div>
  );
}
