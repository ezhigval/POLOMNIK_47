import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { LegalDocumentView } from "@/components/legal-document-view";
import { fetchLegalDocument } from "@/lib/api/legal";
import { legalDocumentPaths, type LegalDocumentType } from "@/lib/operator-config";

const slugToType: Record<string, LegalDocumentType> = {
  "privacy-policy": "privacy_policy",
  "personal-data-consent": "personal_data",
  "distribution-consent": "distribution",
  "marketing-consent": "marketing",
  "cookie-policy": "cookie",
  terms: "terms",
  offer: "offer",
};

type PageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateStaticParams() {
  return Object.keys(slugToType).map((slug) => ({ slug }));
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;
  const type = slugToType[slug];
  if (!type) {
    return { title: "Документ не найден" };
  }
  try {
    const doc = await fetchLegalDocument(type);
    return {
      title: doc.title,
      alternates: { canonical: legalDocumentPaths[type] },
    };
  } catch {
    return { title: "Юридический документ" };
  }
}

export default async function LegalDocumentPage({ params }: PageProps) {
  const { slug } = await params;
  const type = slugToType[slug];
  if (!type) {
    notFound();
  }

  let doc;
  try {
    doc = await fetchLegalDocument(type);
  } catch {
    notFound();
  }

  return <LegalDocumentView doc={doc} type={type} />;
}
