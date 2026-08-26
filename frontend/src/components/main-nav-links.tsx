"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { isMainNavLinkActive, mainNavLinks } from "@/lib/site-nav";

export function MainNavLinks() {
  const pathname = usePathname();

  return (
    <>
      {mainNavLinks.map((link) => {
        const active = isMainNavLinkActive(link.href, pathname);
        return (
          <Link
            key={link.href}
            href={link.href}
            aria-current={active ? "page" : undefined}
            className={`nav-link ${active ? "bg-brand-50 font-medium text-brand-900" : ""}`}
          >
            {link.label}
          </Link>
        );
      })}
    </>
  );
}
