import Link from "next/link";

export function HomeSpaceSwitcher({ guide }: { guide: boolean }) {
  return (
    <nav aria-label="Режим главной страницы" className="mx-auto flex max-w-6xl justify-center px-4 pt-5 sm:justify-start">
      <div className="inline-flex rounded-full border border-stone-200 bg-stone-100 p-1">
        {[
          { href: "/", label: "Подобрать поездку", selected: !guide },
          { href: "/?mode=guide", label: "Путеводитель", selected: guide },
        ].map(({ href, label, selected }) => (
          <Link key={href} href={href} scroll={false} aria-current={selected ? "page" : undefined}
            className={`rounded-full px-4 py-2.5 text-sm font-semibold focus-visible:outline-2 sm:px-6 ${selected ? "bg-brand-900 text-white" : "text-stone-600 hover:bg-white"}`}>
            {label}
          </Link>
        ))}
      </div>
    </nav>
  );
}
