export function portraitUrl(value?: string): string {
  const src = value?.trim() ?? "";
  if (src.startsWith("/") && !src.startsWith("//") && !src.includes("\\")) return src;
  try {
    return ["https:", "http:"].includes(new URL(src).protocol) ? src : "";
  } catch {
    return "";
  }
}
