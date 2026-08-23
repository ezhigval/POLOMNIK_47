export type CookieConsentChoice = "all" | "essential" | "reject";

export const COOKIE_CONSENT_KEY = "palomnik_cookie_consent";
export const COOKIE_CONSENT_EVENT = "palomnik:cookie-consent-open";
export const COOKIE_CONSENT_CHANGED_EVENT = "palomnik:cookie-consent-changed";

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

export function subscribeCookieConsent(onStoreChange: () => void): () => void {
  if (typeof window === "undefined") {
    return () => {};
  }
  window.addEventListener(COOKIE_CONSENT_CHANGED_EVENT, onStoreChange);
  return () => window.removeEventListener(COOKIE_CONSENT_CHANGED_EVENT, onStoreChange);
}

export function setCookieConsent(choice: CookieConsentChoice): void {
  window.localStorage.setItem(COOKIE_CONSENT_KEY, choice);
  window.dispatchEvent(new CustomEvent(COOKIE_CONSENT_CHANGED_EVENT, { detail: choice }));
}

export function allowsAnalyticsCookies(choice: CookieConsentChoice | null): boolean {
  return choice === "all";
}
