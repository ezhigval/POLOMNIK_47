import { apiGet } from "./client";

export type LegalDocumentSummary = {
  id: string;
  type: string;
  version: string;
  title: string;
  published_at: string;
  updated_at: string;
  is_active: boolean;
};

export type LegalDocument = LegalDocumentSummary & {
  content: string;
};

export async function fetchLegalDocument(type: string): Promise<LegalDocument> {
  return apiGet<LegalDocument>(`/legal/documents/${type}`);
}

export async function fetchLegalDocuments(): Promise<LegalDocumentSummary[]> {
  return apiGet<LegalDocumentSummary[]>("/legal/documents");
}

export async function recordConsent(payload: {
  consent_type: string;
  request_id?: string;
}): Promise<void> {
  await fetch("/api/consents", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
}
