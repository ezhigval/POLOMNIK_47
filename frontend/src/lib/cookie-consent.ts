export type CookieConsentChoice = "all" | "essential" | "reject";

export const COOKIE_CONSENT_KEY = "palomnik_cookie_consent";
export const COOKIE_CONSENT_EVENT = "palomnik:cookie-consent-open";

export function getCookieConsent(): CookieConsentChoice | null {
  if (typeof window === "undefined") {
    return null;
  }
  const raw = window.localStorage.getItem(COOKIE_CONSENT_KEY);
  if (raw === "all" || raw === "essential" || raw === "reject") {
    return raw;
  }
  return null;
}

export function setCookieConsent(choice: CookieConsentChoice): void {
  window.localStorage.setItem(COOKIE_CONSENT_KEY, choice);
  window.dispatchEvent(new CustomEvent("palomnik:cookie-consent-changed", { detail: choice }));
}

export function allowsAnalyticsCookies(choice: CookieConsentChoice | null): boolean {
  return choice === "all";
}
