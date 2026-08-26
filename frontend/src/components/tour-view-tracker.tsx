"use client";

import { useEffect } from "react";
import { trackTourView } from "@/lib/analytics";

type TourViewTrackerProps = {
  tourId: string;
  title: string;
};

export function TourViewTracker({ tourId, title }: TourViewTrackerProps) {
  useEffect(() => {
    trackTourView(tourId, title);
  }, [tourId, title]);

  useEffect(() => {
    void fetch("/api/viewed-tours", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ tour_id: tourId }),
    });
  }, [tourId]);

  return null;
}
