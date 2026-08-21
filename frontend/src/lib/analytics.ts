export const analyticsConfig = {
  yandexMetrikaId: process.env.NEXT_PUBLIC_YM_ID?.trim() || "",
  googleAnalyticsId: process.env.NEXT_PUBLIC_GA_ID?.trim() || "",
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

export function trackBookingSubmit(tourId: string, peopleCount: number): void {
  trackEvent("booking_submit", { tour_id: tourId, people_count: peopleCount });
}
