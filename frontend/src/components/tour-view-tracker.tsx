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

  return null;
}
