import Link from "next/link";
import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { PublishLegalDocumentForm } from "@/components/management/publish-legal-document-form";
import { StatusBadge } from "@/components/management/status-badge";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";
import { listManagementConsents, listManagementLegalDocuments } from "@/lib/api/management";
import { legalDocumentPaths } from "@/lib/operator-config";

export default async function ManagementLegalPage() {
  if (!(await canAccessManagementPage([PERM.content]))) {
    return <ManagementNoAccess />;
  }

  const documents = await listManagementLegalDocuments();
  const consents = await listManagementConsents(1, 30);

  return (
    <div className="space-y-8">
      <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
        <ManagementPanel title="Версии документов" description={`${documents.length} записей`}>
          <ManagementTable>
            <ManagementTableHead>
              <ManagementTh>Документ</ManagementTh>
              <ManagementTh>Версия</ManagementTh>
              <ManagementTh>Статус</ManagementTh>
              <ManagementTh />
            </ManagementTableHead>
            <tbody>
              {documents.length === 0 ? (
                <ManagementEmptyRow colSpan={4}>Документов пока нет.</ManagementEmptyRow>
              ) : (
                documents.map((doc) => (
                  <tr key={doc.id} className="border-b border-stone-100 align-top last:border-0">
                    <td className="px-4 py-4">
                      <p className="font-medium text-stone-900">{doc.title}</p>
                      <p className="mt-1 text-xs text-stone-500">{doc.type}</p>
                    </td>
                    <td className="px-4 py-4 text-sm">{doc.version}</td>
                    <td className="px-4 py-4">
                      <StatusBadge variant={doc.is_active ? "success" : "warning"}>
                        {doc.is_active ? "Активна" : "Архив"}
                      </StatusBadge>
                    </td>
                    <td className="px-4 py-4 text-right text-sm">
                      {doc.is_active && legalDocumentPaths[doc.type as keyof typeof legalDocumentPaths] ? (
                        <Link
                          href={legalDocumentPaths[doc.type as keyof typeof legalDocumentPaths]}
                          className="text-brand-800 underline underline-offset-2"
                          target="_blank"
                        >
                          Открыть
                        </Link>
                      ) : null}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </ManagementTable>
        </ManagementPanel>

        <PublishLegalDocumentForm />
      </div>

      <ManagementPanel title="Факты согласий" description={`Показаны последние ${consents.data.length}`}>
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Тип</ManagementTh>
            <ManagementTh>Версия</ManagementTh>
            <ManagementTh>Дата</ManagementTh>
            <ManagementTh>User / Request</ManagementTh>
          </ManagementTableHead>
          <tbody>
            {consents.data.length === 0 ? (
              <ManagementEmptyRow colSpan={4}>Согласий пока нет.</ManagementEmptyRow>
            ) : (
              consents.data.map((item) => (
                <tr key={item.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-3 text-sm font-medium">{item.consent_type}</td>
                  <td className="px-4 py-3 text-sm">{item.document_version}</td>
                  <td className="px-4 py-3 text-sm text-stone-600">
                    {new Date(item.accepted_at).toLocaleString("ru-RU")}
                  </td>
                  <td className="px-4 py-3 font-mono text-xs text-stone-500">
                    {item.user_id?.slice(0, 8) ?? "—"} / {item.request_id?.slice(0, 8) ?? "—"}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>
    </div>
  );
}
