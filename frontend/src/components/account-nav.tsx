"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { accountNavLinks } from "@/lib/site-nav";

export function AccountNav() {
  const pathname = usePathname();

  return (
    <nav
      className="flex flex-wrap gap-2 border-b border-stone-200 pb-4"
      aria-label="Разделы личного кабинета"
    >
      {accountNavLinks.map((link) => {
        const active =
          link.href === "/account"
            ? pathname === "/account"
            : pathname === link.href || pathname.startsWith(`${link.href}/`);
        return (
          <Link
            key={link.href}
            href={link.href}
            aria-current={active ? "page" : undefined}
            className={`rounded-full px-4 py-2 text-sm font-medium transition ${
              active
                ? "bg-brand-800 text-white"
                : "bg-white text-stone-700 ring-1 ring-stone-200 hover:bg-stone-50"
            }`}
          >
            {link.label}
          </Link>
        );
      })}
    </nav>
  );
}
