/**
 * Canonical social OAuth paths for owner console setup.
 * Redirect URIs are on the public site (Next.js), not api.* subdomain.
 * @see docs/OAUTH_SETUP.md
 */
export const SOCIAL_AUTH_PATHS = {
  yandexCallback: "/api/auth/social/yandex/callback",
  vkCallback: "/api/auth/social/vk/callback",
  maxCallback: "/api/auth/social/max/callback",
  telegram: "/api/auth/social/telegram",
} as const;

export const SOCIAL_AUTH_PUBLIC_ORIGIN = "https://tikhvin-palomnik.ru";

export function socialRedirectURI(
  key: keyof typeof SOCIAL_AUTH_PATHS,
  origin: string = SOCIAL_AUTH_PUBLIC_ORIGIN,
): string {
  return `${origin.replace(/\/$/, "")}${SOCIAL_AUTH_PATHS[key]}`;
}
