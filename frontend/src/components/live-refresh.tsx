"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

type LiveRefreshProps = {
  intervalMs?: number;
};

export function LiveRefresh({ intervalMs = 3000 }: LiveRefreshProps) {
  const router = useRouter();

  useEffect(() => {
    if (process.env.NEXT_PUBLIC_LIVE_REFRESH !== "1") {
      return;
    }

    const id = window.setInterval(() => {
      router.refresh();
    }, intervalMs);

    return () => window.clearInterval(id);
  }, [intervalMs, router]);

  return null;
}
