import Link from "next/link";
import { notFound, redirect } from "next/navigation";
import type { Metadata } from "next";
import { AdminLoginForm } from "@/components/management/admin-login-form";
import { isManagementConfigured } from "@/lib/api/management";
import { isAdminAuthenticated } from "@/lib/auth/admin-session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Вход в админку",
  robots: { index: false, follow: false },
};

export default async function ManagementLoginPage() {
  if (!isManagementConfigured()) {
    notFound();
  }

  if (await isAdminAuthenticated()) {
    redirect("/management");
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-stone-100 px-4 py-12">
      <Link href="/" className="mb-8 text-sm text-stone-500 hover:text-brand-800">
        ← На сайт
      </Link>
      <AdminLoginForm />
    </div>
  );
}
