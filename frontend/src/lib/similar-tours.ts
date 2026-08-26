import type { Tour } from "@/lib/api/tours";
import type { TourFilterValues } from "@/lib/tour-filters";

const MS_PER_DAY = 24 * 60 * 60 * 1000;

export function rankSimilarTours(tours: Tour[], filters: TourFilterValues, limit = 6): Tour[] {
  const people = Number.parseInt(filters.people || filters.min_slots || "", 10);
  const dateFrom = utcDay(filters.date_from);
  const dateTo = utcDay(filters.date_to);

  return [...tours]
    .map((tour) => ({
      tour,
      dateGap: dateGapDays(tour, dateFrom, dateTo),
      slotGap: slotGap(tour.slots_left, people),
      start: utcDay(tour.date_start) ?? Number.MAX_SAFE_INTEGER,
    }))
    .sort((a, b) => a.dateGap - b.dateGap || a.slotGap - b.slotGap || a.start - b.start)
    .slice(0, limit)
    .map((item) => item.tour);
}

function slotGap(slotsLeft: number, people: number): number {
  if (!Number.isFinite(people) || people <= 0) {
    return 0;
  }
  return Math.max(0, people - slotsLeft);
}

function dateGapDays(tour: Tour, windowStart: number | null, windowEnd: number | null): number {
  const tourStart = utcDay(tour.date_start);
  const tourEnd = utcDay(tour.date_end);
  if (tourStart == null || tourEnd == null) {
    return Number.MAX_SAFE_INTEGER;
  }
  if (windowStart == null && windowEnd == null) {
    return 0;
  }

  const start = windowStart ?? Number.NEGATIVE_INFINITY;
  const end = windowEnd ?? Number.POSITIVE_INFINITY;
  if (tourEnd >= start && tourStart <= end) {
    return 0;
  }
  if (tourEnd < start) {
    return Math.round((start - tourEnd) / MS_PER_DAY);
  }
  return Math.round((tourStart - end) / MS_PER_DAY);
}

function utcDay(value?: string | null): number | null {
  if (!value) {
    return null;
  }
  const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(value);
  if (!match) {
    return null;
  }
  return Date.UTC(Number(match[1]), Number(match[2]) - 1, Number(match[3]));
}
