"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import type { MouseEvent } from "react";
import { canUseHistoryBack } from "@/lib/navigation-history";

type SmartBackButtonProps = {
  fallbackHref: string;
  label?: string;
  className?: string;
};

export function SmartBackButton({
  fallbackHref,
  label = "Назад",
  className = "inline-flex items-center gap-1 text-sm text-stone-500 transition hover:text-brand-800",
}: SmartBackButtonProps) {
  const router = useRouter();
  const pathname = usePathname();

  function onClick(event: MouseEvent<HTMLAnchorElement>) {
    if (canUseHistoryBack(pathname)) {
      event.preventDefault();
      router.back();
    }
  }

  return (
    <Link href={fallbackHref} onClick={onClick} className={className} aria-label={label}>
      <span aria-hidden="true">←</span>
      {label}
    </Link>
  );
}
