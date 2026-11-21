import "server-only";
import { cache } from "react";
import { apiUrl, requestJson, type DataEnvelope } from "./api/client";
import { getApiBaseUrl } from "./api/base-url";
import { defaultLogo, isUploadedLogo } from "./brand";

export const getSiteLogo = cache(async (): Promise<string> => {
  try {
    const { data } = await requestJson<DataEnvelope<{logo_path?: string}>>(apiUrl("/site-settings"), { cache: "no-store" });
    return data.logo_path && isUploadedLogo(data.logo_path) ? data.logo_path : defaultLogo;
  } catch {
    return defaultLogo;
  }
});

export function fetchUploadedLogo(path: string) {
  if (!isUploadedLogo(path)) throw new Error("Недопустимый путь к логотипу");
  return fetch(new URL(path, getApiBaseUrl()), { redirect: "error", cache: "no-store" });
}
