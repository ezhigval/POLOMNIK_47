import { redirect } from "next/navigation";
import { BootstrapHomeButton } from "@/components/management/cms/bootstrap-home-button";
import { listManagementCmsPagesOrEmpty } from "@/lib/api/management";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

export default async function ManagementContentPage() {
  if (!(await canAccessManagementPage([PERM.content]))) {
    return <ManagementNoAccess />;
  }
  const { pages, unavailable } = await listManagementCmsPagesOrEmpty();

  if (unavailable) {
    return (
      <p className="rounded-2xl border border-stone-200 bg-white p-5 text-sm text-stone-600">
        Редактор главной сейчас недоступен. Обновите страницу или проверьте API.
      </p>
    );
  }

  const home = pages.find((page) => page.slug === "home");
  if (home) {
    redirect(`/management/content/${home.id}`);
  }

  return (
    <div className="max-w-xl space-y-4">
      <div>
        <h2 className="text-xl font-semibold text-stone-900">Главная</h2>
        <p className="mt-1 text-sm text-stone-600">
          Страницы главной ещё нет. Создайте её один раз — дальше только блоки, тексты и SEO, без
          новых страниц.
        </p>
      </div>
      <BootstrapHomeButton />
    </div>
  );
}
