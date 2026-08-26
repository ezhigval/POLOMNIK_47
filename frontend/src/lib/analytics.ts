export const analyticsConfig = {
  yandexMetrikaId: process.env.NEXT_PUBLIC_YM_ID?.trim() || "",
  googleAnalyticsId: process.env.NEXT_PUBLIC_GA_ID?.trim() || "",
  /** Webvisor — только при явном NEXT_PUBLIC_YM_WEBVISOR=1 */
  yandexWebvisor: process.env.NEXT_PUBLIC_YM_WEBVISOR === "1",
  /** Карта кликов включена вместе с Метрикой; NEXT_PUBLIC_YM_CLICKMAP=0 выключает. */
  yandexClickmap: process.env.NEXT_PUBLIC_YM_CLICKMAP !== "0",
};

export function isAnalyticsEnabled(): boolean {
  return Boolean(analyticsConfig.yandexMetrikaId || analyticsConfig.googleAnalyticsId);
}

type AnalyticsParams = Record<string, string | number | boolean | undefined>;

declare global {
  interface Window {
    ym?: (id: number, method: string, goal: string, params?: AnalyticsParams) => void;
    gtag?: (...args: unknown[]) => void;
  }
}

export function trackEvent(name: string, params: AnalyticsParams = {}): void {
  if (typeof window === "undefined") {
    return;
  }

  const ymId = Number(analyticsConfig.yandexMetrikaId);
  if (analyticsConfig.yandexMetrikaId && Number.isFinite(ymId)) {
    window.ym?.(ymId, "reachGoal", name, params);
  }

  if (analyticsConfig.googleAnalyticsId) {
    window.gtag?.("event", name, params);
  }
}

export function trackTourView(tourId: string, title: string): void {
  trackEvent("tour_view", { tour_id: tourId, tour_title: title });
}

/** Начало заполнения формы заявки (первый фокус). */
export function trackBeginCheckout(tourId: string): void {
  trackEvent("begin_checkout", { tour_id: tourId });
}

export function trackBookingSubmit(tourId: string, peopleCount: number): void {
  trackEvent("booking_submit", { tour_id: tourId, people_count: peopleCount });
}

/** Клик по телефону / email / чату поддержки. */
export function trackSupportContact(channel: "phone" | "email" | "chat"): void {
  trackEvent("support_contact", { channel });
}
