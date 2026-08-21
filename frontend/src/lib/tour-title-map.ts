import type { ManagementTour } from "@/lib/api/management";

export function buildTourTitleMap(tours: ManagementTour[]): Map<string, string> {
  return new Map(tours.map((tour) => [tour.id, tour.title]));
}

export function tourTitle(map: Map<string, string>, tourId: string): string {
  return map.get(tourId) ?? tourId;
}
