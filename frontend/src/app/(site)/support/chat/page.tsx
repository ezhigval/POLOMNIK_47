import type { Metadata } from "next";
import { PageIntro } from "@/components/page-intro";
import { SupportChat } from "@/components/support-chat";
import { fetchSupportThread } from "@/lib/auth/user-features";
import { getSessionUser } from "@/lib/auth/session";

export const metadata: Metadata = {
  title: "Чат поддержки",
  robots: { index: false, follow: false },
};

export default async function SupportChatPage() {
  const user = await getSessionUser();
  let thread = null;

  if (user) {
    try {
      thread = await fetchSupportThread();
    } catch {
      thread = null;
    }
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6 px-4 py-8 sm:py-12">
      <PageIntro
        backHref="/support"
        backLabel="← Справочник поддержки"
        title="Чат поддержки"
        description="Задайте вопрос — менеджер ответит в рабочее время. История переписки сохраняется в аккаунте."
      />
      <SupportChat initialThread={thread} isAuthenticated={Boolean(user)} />
    </div>
  );
}
