import { BrandMark } from "@/components/brand-mark";
import Link from "next/link";
import { redirect } from "next/navigation";
import { ManagementNav } from "@/components/management/management-nav";
import { AdminLogoutButton } from "@/components/management/admin-logout-button";
import { getManagementSession, isManagementConfigured } from "@/lib/api/management";
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

  const session = await getManagementSession().catch(() => ({
    full_admin: false,
    permissions: [] as string[],
  }));

  return (
    <div className="min-h-full bg-stone-50">
      <header className="border-b border-stone-200 bg-white shadow-sm">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
          <Link href="/management" className="flex items-center gap-2 max-w-[14rem] text-sm font-semibold leading-snug text-stone-900 hover:text-brand-800 sm:max-w-none">
            <BrandMark className="size-12" />
            <span>{siteConfig.name}<span className="mt-0.5 block text-xs font-normal text-stone-500">Панель управления</span></span>
          </Link>
          <div className="flex items-center gap-3">
            <Link href="/" target="_blank" className="hidden rounded-xl border border-stone-200 px-3 py-2 text-sm font-medium text-stone-600 hover:bg-stone-50 sm:inline-flex">
              Открыть сайт ↗
            </Link>
            <AdminLogoutButton />
          </div>
        </div>
        <div className="border-t border-stone-100">
          <div className="mx-auto max-w-7xl px-4 sm:px-6">
            <ManagementNav fullAdmin={session.full_admin} permissions={session.permissions ?? []} />
          </div>
        </div>
      </header>

      <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-8">
        <div className="mb-6">
          <p className="text-xs font-medium uppercase tracking-widest text-stone-500">Управление сайтом и службой</p>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight text-stone-900">Управление</h1>
        </div>
        {children}
      </div>
    </div>
  );
}
