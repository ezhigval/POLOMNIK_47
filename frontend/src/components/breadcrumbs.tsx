import Link from "next/link";
import { JsonLd } from "@/components/structured-data";
import { absoluteUrl } from "@/lib/site-config";

export type BreadcrumbItem = {
  name: string;
  href?: string;
};

export function Breadcrumbs({ items, className = "mb-6" }: { items: BreadcrumbItem[]; className?: string }) {
  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: items.map((item, index) => ({
      "@type": "ListItem",
      position: index + 1,
      name: item.name,
      ...(item.href ? { item: absoluteUrl(item.href) } : {}),
    })),
  };

  return (
    <>
      <JsonLd data={jsonLd} />
      <nav className={`${className} text-sm text-stone-500`} aria-label="Хлебные крошки">
        {items.map((item, index) => {
          const isLast = index === items.length - 1;
          return (
            <span key={`${item.name}-${index}`}>
              {index > 0 ? <span className="mx-2">/</span> : null}
              {item.href && !isLast ? (
                <Link href={item.href} className="hover:text-brand-800">
                  {item.name}
                </Link>
              ) : (
                <span className={isLast ? "text-stone-800" : undefined}>{item.name}</span>
              )}
            </span>
          );
        })}
      </nav>
    </>
  );
}
