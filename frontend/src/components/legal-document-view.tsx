import Link from "next/link";
import type { LegalDocumentType } from "@/lib/operator-config";
import { legalDocumentPaths } from "@/lib/operator-config";
import { DetailPageNav } from "@/components/detail-page-nav";

type LegalDocumentViewProps = {
  doc: {
    title: string;
    version: string;
    published_at: string;
    updated_at: string;
    content: string;
  };
  type: LegalDocumentType;
};

function formatDate(iso: string): string {
  try {
    return new Intl.DateTimeFormat("ru-RU", { dateStyle: "long" }).format(new Date(iso));
  } catch {
    return iso;
  }
}

export function LegalDocumentView({ doc, type }: LegalDocumentViewProps) {
  return (
    <div className="mx-auto max-w-3xl px-4 py-10 sm:py-14">
      <DetailPageNav
        fallbackHref="/legal"
        backLabel="Назад"
        breadcrumbs={[
          { name: "Главная", href: "/" },
          { name: "Юридические документы", href: "/legal" },
          { name: doc.title },
        ]}
      />

      <header className="mb-8 border-b border-stone-200 pb-6">
        <h1 className="font-display text-2xl font-semibold text-stone-900 sm:text-3xl">{doc.title}</h1>
        <dl className="mt-4 flex flex-wrap gap-x-6 gap-y-1 text-xs text-stone-500">
          <div>
            <dt className="inline">Версия: </dt>
            <dd className="inline font-medium text-stone-700">{doc.version}</dd>
          </div>
          <div>
            <dt className="inline">Опубликовано: </dt>
            <dd className="inline">{formatDate(doc.published_at)}</dd>
          </div>
          <div>
            <dt className="inline">Изменено: </dt>
            <dd className="inline">{formatDate(doc.updated_at)}</dd>
          </div>
        </dl>
        <p className="mt-3 text-xs text-stone-500">
          Документ является проектом для проверки оператором/юристом перед публикацией.
        </p>
      </header>

      <article
        className="prose prose-stone max-w-none text-sm leading-7 text-stone-700 prose-headings:font-display prose-a:text-brand-800"
        dangerouslySetInnerHTML={{ __html: doc.content }}
      />

      <footer className="mt-10 border-t border-stone-200 pt-6 text-sm text-stone-500">
        <div className="flex flex-wrap gap-4">
          <Link href={legalDocumentPaths[type]} className="font-medium text-brand-800 hover:underline">
            Постоянная ссылка на документ
          </Link>
          <a
            href={`/api/v1/legal/documents/${type}/download`}
            className="font-medium text-brand-800 hover:underline"
          >
            Скачать (HTML)
          </a>
        </div>
      </footer>
    </div>
  );
}
