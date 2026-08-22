import Link from "next/link";
import { redirect } from "next/navigation";
import { ManagementNav } from "@/components/management/management-nav";
import { AdminLogoutButton } from "@/components/management/admin-logout-button";
import { isManagementConfigured } from "@/lib/api/management";
import { isAdminAuthenticated } from "@/lib/auth/admin-session";
import { siteConfig } from "@/lib/site-config";

export const dynamic = "force-dynamic";

export default async function ProtectedManagementLayout({ children }: { children: React.ReactNode }) {
  if (!isManagementConfigured()) {
    redirect("/");
  }

  if (!(await isAdminAuthenticated())) {
    redirect("/management/login");
  }

  return (
    <div className="min-h-full bg-stone-100">
      <header className="border-b border-stone-200 bg-white">
        <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-x-4 gap-y-2 px-4 py-3">
          <Link href="/management" className="text-sm font-semibold text-stone-900 hover:text-brand-800">
            {siteConfig.name} · Админка
          </Link>
          <div className="flex min-w-0 flex-1 flex-wrap items-center justify-end gap-3">
            <ManagementNav />
            <AdminLogoutButton />
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-6xl px-4 py-8 sm:py-10">
        <div className="mb-8">
          <p className="text-xs font-medium uppercase tracking-widest text-stone-500">Внутренний интерфейс</p>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight text-stone-900">Управление</h1>
        </div>
        {children}
      </div>
    </div>
  );
}
