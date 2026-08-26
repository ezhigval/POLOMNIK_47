"use client";

import { usePathname } from "next/navigation";
import { useEffect } from "react";
import { NAV_CURRENT_PATH_KEY, NAV_PREV_PATH_KEY } from "@/lib/navigation-history";

/** Keeps the previous in-app path in sessionStorage for smart back on detail pages. */
export function NavigationTracker() {
  const pathname = usePathname();

  useEffect(() => {
    const current = sessionStorage.getItem(NAV_CURRENT_PATH_KEY);
    if (current && current !== pathname) {
      sessionStorage.setItem(NAV_PREV_PATH_KEY, current);
    }
    sessionStorage.setItem(NAV_CURRENT_PATH_KEY, pathname);
  }, [pathname]);

  return null;
}
