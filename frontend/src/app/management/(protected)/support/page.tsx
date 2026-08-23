import Link from "next/link";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { listManagementSupportThreads } from "@/lib/api/management";
import { formatDateTime } from "@/lib/format";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

export default async function ManagementSupportPage() {
  if (!(await canAccessManagementPage([PERM.support]))) {
    return <ManagementNoAccess />;
  }
  const threads = await listManagementSupportThreads();

  return (
    <ManagementPanel
      title="Поддержка"
      description={`Всего ${threads.length}. Отвечайте в диалоге — паломник увидит ответ в чате.`}
    >
      <ManagementTable>
        <ManagementTableHead>
          <ManagementTh>Тема</ManagementTh>
          <ManagementTh>Пользователь</ManagementTh>
          <ManagementTh>Статус</ManagementTh>
          <ManagementTh>Обновлён</ManagementTh>
          <ManagementTh />
        </ManagementTableHead>
        <tbody>
          {threads.length === 0 ? (
            <ManagementEmptyRow colSpan={5}>Обращений пока нет.</ManagementEmptyRow>
          ) : (
            threads.map((thread) => (
              <tr key={thread.id} className="border-b border-stone-100 align-top last:border-0">
                <td className="px-4 py-4">
                  <div className="font-medium text-stone-900">{thread.subject}</div>
                  <div className="font-mono text-xs text-stone-400">{thread.id}</div>
                </td>
                <td className="px-4 py-4 font-mono text-xs text-stone-600">{thread.user_id}</td>
                <td className="px-4 py-4">
                  <StatusBadge variant={thread.status === "open" ? "warning" : "neutral"}>
                    {thread.status}
                  </StatusBadge>
                </td>
                <td className="px-4 py-4 whitespace-nowrap text-stone-600">
                  {formatDateTime(thread.updated_at)}
                </td>
                <td className="px-4 py-4 text-right">
                  <Link href={`/management/support/${thread.id}`} className="btn-secondary">
                    Открыть
                  </Link>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </ManagementTable>
    </ManagementPanel>
  );
}
