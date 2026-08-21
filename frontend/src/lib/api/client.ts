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

type DataEnvelope<T> = {
  data: T;
};

type ListEnvelope<T> = {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
    has_next: boolean;
  };
};

async function parseError(response: Response): Promise<ApiError> {
  try {
    const body = (await response.json()) as ErrorEnvelope;
    return new ApiError(
      response.status,
      body.error?.code ?? "UNKNOWN_ERROR",
      body.error?.message ?? "Request failed",
    );
  } catch {
    return new ApiError(response.status, "UNKNOWN_ERROR", response.statusText);
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const response = await fetch(`${getApiBaseUrl()}${path}`, {
    next: { revalidate: 60 },
  });

  if (!response.ok) {
    throw await parseError(response);
  }

  const body = (await response.json()) as DataEnvelope<T>;
  return body.data;
}

export async function apiGetList<T>(path: string): Promise<ListEnvelope<T>> {
  const response = await fetch(`${getApiBaseUrl()}${path}`, {
    next: { revalidate: 60 },
  });

  if (!response.ok) {
    throw await parseError(response);
  }

  return (await response.json()) as ListEnvelope<T>;
}

export async function apiPost<TResponse, TBody>(
  path: string,
  body: TBody,
): Promise<TResponse> {
  const response = await fetch(`${getApiBaseUrl()}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw await parseError(response);
  }

  const payload = (await response.json()) as DataEnvelope<TResponse>;
  return payload.data;
}
