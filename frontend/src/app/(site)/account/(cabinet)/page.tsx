import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { fetchCurrentUser } from "@/lib/api/auth";
import { getAuthToken } from "@/lib/auth/session";
import { userContactLine } from "@/lib/site-nav";

export const metadata: Metadata = {
  title: "Профиль",
};

export default async function AccountProfilePage() {
  const token = await getAuthToken();
  if (!token) {
    redirect("/account/login?returnUrl=%2Faccount");
  }

  let user;
  try {
    user = await fetchCurrentUser(token);
  } catch {
    redirect("/account/login?returnUrl=%2Faccount");
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900">Профиль</h1>
        <p className="mt-2 text-sm text-stone-600">Данные аккаунта. UUID нужен для назначения роли в админке.</p>
      </div>

      <dl className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-500">Имя</dt>
          <dd className="mt-1 text-stone-900">{user.name || "—"}</dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-500">Контакт</dt>
          <dd className="mt-1 text-stone-900">{userContactLine(user)}</dd>
        </div>
        <div>
          <dt className="text-xs uppercase tracking-wide text-stone-500">UUID пользователя</dt>
          <dd className="mt-1 break-all font-mono text-sm text-stone-900">{user.id}</dd>
        </div>
      </dl>
    </div>
  );
}
