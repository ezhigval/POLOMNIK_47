import Link from "next/link";
import type { Metadata } from "next";
import { SectionHeading } from "@/components/section-heading";
import { fetchLegalDocuments } from "@/lib/api/legal";
import { legalDocumentPaths } from "@/lib/operator-config";

export const metadata: Metadata = {
  title: "Юридические документы",
  description: "Политика обработки персональных данных, согласия и cookie policy.",
  alternates: { canonical: "/legal" },
};

export default async function LegalIndexPage() {
  let documents: Awaited<ReturnType<typeof fetchLegalDocuments>> = [];
  try {
    documents = await fetchLegalDocuments();
  } catch {
    documents = [];
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-10 sm:py-14">
      <div className="mb-6">
        <Link href="/" className="text-sm text-stone-500 transition hover:text-brand-800">
          ← На главную
        </Link>
      </div>

      <SectionHeading
        title="Юридические документы"
        description="Политики и согласия в отношении персональных данных. Актуальные версии с указанием даты публикации."
      />

      {documents.length === 0 ? (
        <p className="mt-8 text-sm text-stone-600">Документы загружаются… Если список пуст, обратитесь в поддержку.</p>
      ) : (
        <ul className="mt-8 divide-y divide-stone-200 rounded-2xl border border-stone-200 bg-white">
          {documents.map((doc) => (
            <li key={doc.id}>
              <Link
                href={legalDocumentPaths[doc.type as keyof typeof legalDocumentPaths] ?? `/legal/${doc.type}`}
                className="flex flex-col gap-1 px-5 py-4 transition hover:bg-stone-50 sm:flex-row sm:items-center sm:justify-between"
              >
                <span className="font-medium text-stone-900">{doc.title}</span>
                <span className="text-xs text-stone-500">v{doc.version}</span>
              </Link>
            </li>
          ))}
        </ul>
      )}

      <p className="mt-8 text-xs leading-5 text-stone-500">
        Документы являются проектом для проверки оператором/юристом. Не заменяют юридическое заключение.
      </p>
    </div>
  );
}
