import Link from "next/link";
import { notFound } from "next/navigation";
import { ManagementPanel } from "@/components/management/management-panel";
import { ReplySupportForm } from "@/components/management/reply-support-form";
import { SupportDraftBox } from "@/components/management/support-draft-box";
import { StatusBadge } from "@/components/management/status-badge";
import { ApiError } from "@/lib/api/client";
import { getManagementSupportThread } from "@/lib/api/management";
import { formatDateTime } from "@/lib/format";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

type PageProps = {
  params: Promise<{ id: string }>;
};

export default async function ManagementSupportThreadPage({ params }: PageProps) {
  if (!(await canAccessManagementPage([PERM.support]))) {
    return <ManagementNoAccess />;
  }
  const { id } = await params;

  let thread;
  try {
    thread = await getManagementSupportThread(id);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      notFound();
    }
    throw error;
  }

  const messages = thread.messages ?? [];

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <Link href="/management/support" className="text-sm text-stone-500 hover:text-brand-800">
            ← Все обращения
          </Link>
          <h2 className="mt-2 text-xl font-semibold text-stone-900">{thread.subject}</h2>
          <p className="mt-1 font-mono text-xs text-stone-400">{thread.id}</p>
        </div>
        <StatusBadge variant={thread.status === "open" ? "warning" : "neutral"}>
          {thread.status}
        </StatusBadge>
      </div>

      <ManagementPanel
        title="Переписка"
        description={`Пользователь: ${thread.user_id} · обновлено ${formatDateTime(thread.updated_at)}. Реплай в боте на уведомление с id диалога — то же сообщение.`}
      >
        <div className="space-y-3 px-4 py-4">
          {messages.length === 0 ? (
            <p className="text-sm text-stone-500">Сообщений пока нет.</p>
          ) : (
            messages.map((message) => {
              const staff = message.sender_type === "staff";
              return (
                <div
                  key={message.id}
                  className={`rounded-xl px-4 py-3 text-sm leading-6 ${
                    staff ? "bg-brand-50 text-stone-900" : "bg-stone-50 text-stone-800"
                  }`}
                >
                  <div className="mb-1 flex flex-wrap items-center justify-between gap-2 text-xs text-stone-500">
                    <span className="font-medium">{staff ? "Сотрудник" : "Паломник"}</span>
                    <span>{formatDateTime(message.created_at)}</span>
                  </div>
                  <p className="whitespace-pre-wrap">{message.body}</p>
                </div>
              );
            })
          )}
        </div>
        <div className="space-y-4 border-t border-stone-100 px-4 py-4">
          <SupportDraftBox threadId={thread.id} />
          <ReplySupportForm threadId={thread.id} />
        </div>
      </ManagementPanel>
    </div>
  );
}
