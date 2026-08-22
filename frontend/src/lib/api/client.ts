import { getApiBaseUrl } from "./base-url";

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

type ErrorEnvelope = {
  error: {
    code: string;
    message: string;
  };
};

export type DataEnvelope<T> = {
  data: T;
};

export type ListEnvelope<T> = {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
    has_next: boolean;
  };
};

const defaultCache: RequestInit =
  process.env.NODE_ENV === "development" || process.env.NEXT_PUBLIC_LIVE_REFRESH === "1"
    ? { cache: "no-store" }
    : { next: { revalidate: 60 } };

function isFormDataBody(body: BodyInit | null | undefined): boolean {
  return typeof FormData !== "undefined" && body instanceof FormData;
}

export async function requestJson<T>(url: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body != null && !isFormDataBody(init.body) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(url, {
    ...defaultCache,
    ...init,
    headers,
  });

  if (response.status === 204) {
    return undefined as T;
  }

  const text = await response.text();
  let payload: unknown = null;
  if (text) {
    try {
      payload = JSON.parse(text);
    } catch {
      payload = null;
    }
  }

  if (!response.ok) {
    const errorBody = payload as ErrorEnvelope | null;
    throw new ApiError(
      response.status,
      errorBody?.error?.code ?? "UNKNOWN_ERROR",
      errorBody?.error?.message ?? "Не удалось выполнить запрос",
    );
  }

  return payload as T;
}

export function apiUrl(path: string): string {
  return `${getApiBaseUrl()}${path}`;
}

export async function apiGet<T>(path: string): Promise<T> {
  const body = await requestJson<DataEnvelope<T>>(apiUrl(path));
  return body.data;
}

export async function apiGetList<T>(path: string): Promise<ListEnvelope<T>> {
  return requestJson<ListEnvelope<T>>(apiUrl(path));
}

export async function apiPost<TResponse, TBody>(path: string, body: TBody): Promise<TResponse> {
  const payload = await requestJson<DataEnvelope<TResponse>>(apiUrl(path), {
    method: "POST",
    body: JSON.stringify(body),
  });
  return payload.data;
}
