export function getApiBaseUrl(): string {
  if (typeof window === "undefined") {
    if (process.env.API_INTERNAL_URL) {
      return process.env.API_INTERNAL_URL;
    }
    const pub = process.env.NEXT_PUBLIC_API_URL;
    if (pub && pub.startsWith("http")) {
      return pub;
    }
    return "http://127.0.0.1:8080/api/v1";
  }

  return "/api/v1";
}
