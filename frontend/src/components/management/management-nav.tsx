"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { visibleManagementNav } from "@/lib/management-access";

type Props = {
  fullAdmin?: boolean;
  permissions?: string[];
};

const sections = [
  { title: "Общее", paths: ["/management"] },
  { title: "Контент", paths: ["/management/content", "/management/news", "/management/smm", "/management/reviews", "/management/legal"] },
  { title: "Работа с людьми", paths: ["/management/tours", "/management/bookings", "/management/support"] },
  { title: "Система", paths: ["/management/integrations", "/management/ai", "/management/settings"] },
];

export function ManagementNav({ fullAdmin = false, permissions = [] }: Props) {
  const pathname = usePathname();
  const links = visibleManagementNav({ full_admin: fullAdmin, permissions });

  return (
    <nav className="overflow-x-auto overscroll-x-contain" aria-label="Разделы панели управления">
      <div className="flex min-w-max items-center gap-5 py-2 text-sm">
        {sections.map((section) => {
          const sectionLinks = section.paths.flatMap((path) => links.filter((link) => link.href === path));
          if (sectionLinks.length === 0) return null;
          return (
            <div key={section.title} className="flex items-center gap-1 border-r border-stone-200 pr-5 last:border-r-0 last:pr-0"
              role="group" aria-label={section.title}>
              <span className="mr-1 hidden text-xs font-semibold uppercase tracking-wide text-stone-400 xl:inline">
                {section.title}
              </span>
              {sectionLinks.map((link) => {
                const active = link.exact
                  ? pathname === link.href
                  : pathname === link.href || pathname.startsWith(`${link.href}/`);
                return (
                  <Link key={link.href} href={link.href} aria-current={active ? "page" : undefined}
                    className={`inline-flex min-h-10 items-center whitespace-nowrap rounded-xl px-3 py-2 font-medium transition focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-700 ${active
                      ? "bg-brand-800 text-white shadow-sm"
                      : "text-stone-600 hover:bg-stone-100 hover:text-stone-900"}`}>
                    {link.label}
                  </Link>
                );
              })}
            </div>
          );
        })}
      </div>
    </nav>
  );
}
