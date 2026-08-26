"use client";

import { usePathname } from "next/navigation";
import { Breadcrumbs, type BreadcrumbItem } from "@/components/breadcrumbs";
import { accountNavLinks } from "@/lib/site-nav";

export function AccountBreadcrumbs() {
  const pathname = usePathname();

  const section = accountNavLinks.find((link) =>
    link.href === "/account"
      ? pathname === "/account"
      : pathname === link.href || pathname.startsWith(`${link.href}/`),
  );

  const items: BreadcrumbItem[] = [
    { name: "Главная", href: "/" },
    { name: "Личный кабинет", href: section && section.href !== "/account" ? "/account" : undefined },
  ];

  if (section && section.href !== "/account") {
    items.push({ name: section.label });
  } else if (pathname === "/account") {
    items[1] = { name: "Профиль" };
  }

  return <Breadcrumbs items={items} className="mb-0" />;
}
