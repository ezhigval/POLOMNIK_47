"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/management", label: "Обзор", exact: true },
  { href: "/management/content", label: "Контент" },
  { href: "/management/news", label: "Новости" },
  { href: "/management/tours", label: "Туры" },
  { href: "/management/bookings", label: "Заявки" },
  { href: "/management/reviews", label: "Отзывы" },
  { href: "/management/integrations", label: "Синхронизация" },
  { href: "/management/settings", label: "Настройки" },
];

export function ManagementNav() {
  const pathname = usePathname();

  return (
    <nav className="flex flex-wrap gap-1 text-sm" aria-label="Навигация управления">
      {links.map((link) => {
        const active =
          link.exact ? pathname === link.href : pathname === link.href || pathname.startsWith(`${link.href}/`);

        return (
          <Link
            key={link.href}
            href={link.href}
            className={`rounded-full px-3 py-1.5 transition ${
              active
                ? "bg-brand-800 font-medium text-white"
                : "text-stone-600 hover:bg-stone-100 hover:text-stone-900"
            }`}
            aria-current={active ? "page" : undefined}
          >
            {link.label}
          </Link>
        );
      })}
      <Link
        href="/"
        className="rounded-full px-3 py-1.5 text-stone-500 transition hover:bg-stone-100 hover:text-stone-800"
      >
        На сайт →
      </Link>
    </nav>
  );
}
