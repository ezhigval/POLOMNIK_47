import type { Metadata } from "next";
import { AccountConsentsPanel } from "@/components/account/account-consents-panel";
import { PageIntro } from "@/components/page-intro";

export const metadata: Metadata = {
  title: "Согласия",
  robots: { index: false, follow: false },
};

export default function AccountConsentsPage() {
  return (
    <div className="space-y-6">
      <PageIntro
        title="Согласия и персональные данные"
        description="История согласий, настройки рекламы и прекращение распространения."
      />
      <AccountConsentsPanel />
    </div>
  );
}
