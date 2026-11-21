export const defaultLogo = "/images/service-logo.png";

const uploadedLogo = /^\/uploads\/[0-9]{4}\/[0-9]{2}\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\.png$/;

export function isUploadedLogo(path: string) {
  return uploadedLogo.test(path);
}

export function logoSrc(path?: string) {
  return path && isUploadedLogo(path) ? `/api/brand?file=${encodeURIComponent(path)}` : defaultLogo;
}
