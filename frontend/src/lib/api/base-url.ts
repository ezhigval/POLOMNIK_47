export function getApiBaseUrl(): string {
  if (typeof window === "undefined" && process.env.API_INTERNAL_URL) {
    return process.env.API_INTERNAL_URL;
  }

  return process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";
}
