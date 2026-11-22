import { ManagementCreateDialog } from "@/components/management/management-create-dialog";
import { ManagementNewsTable } from "@/components/management/management-news-table";
import { ManagementPanel } from "@/components/management/management-panel";
import { listManagementNews } from "@/lib/api/management";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

export default async function ManagementNewsPage() {
  if (!(await canAccessManagementPage([PERM.content]))) {
    return <ManagementNoAccess />;
  }
  const articles = await listManagementNews();

  return (
    <div>
      <ManagementPanel
        title="Новости"
        description={`${articles.length} статей. Звёздочка закрепляет на главной и в ленте (не больше трёх: одна главная и две рядом).`}
        action={<ManagementCreateDialog kind="news" />}
      >
        <ManagementNewsTable articles={articles} />
      </ManagementPanel>
    </div>
  );
}
