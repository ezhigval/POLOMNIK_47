import Link from "next/link";
import { featuredRoute } from "@/lib/featured-route";
import type { FeaturedRouteBlockContent } from "@/lib/api/cms";

type FeaturedRouteSectionProps = {
  content?: FeaturedRouteBlockContent;
};

function featuredCatalogHref(cmsHref?: string): string {
  const href = cmsHref?.trim();
  if (!href || href === "/search" || href.includes("destination=tikhvin")) {
    return featuredRoute.ctaHref;
  }
  return href;
}

export function FeaturedRouteSection({ content }: FeaturedRouteSectionProps = {}) {
  const route = {
    ...featuredRoute,
    ...content,
    days: content?.days?.length ? content.days : featuredRoute.days,
    ctaHref: featuredCatalogHref(content?.ctaHref),
  };
  return (
    <section
      id={featuredRoute.id}
      className="scroll-mt-24 overflow-hidden rounded-3xl border border-amber-200/80 bg-white shadow-sm ring-1 ring-amber-100"
    >
      <div className="grid md:grid-cols-[1.15fr_0.85fr]">
        <div className="bg-gradient-to-br from-brand-950 via-brand-900 to-brand-800 px-6 py-8 text-white sm:px-8 sm:py-10">
          <p className="inline-flex items-center gap-2 rounded-full border border-amber-200/30 bg-white/10 px-3 py-1 text-xs font-medium uppercase tracking-widest text-amber-100">
            <span className="size-1.5 rounded-full bg-amber-300" />
            {route.eyebrow}
          </p>
          <h2 className="mt-4 font-display text-3xl font-semibold tracking-tight sm:text-4xl">
            {route.title}
          </h2>
          <p className="mt-3 text-sm text-brand-100">
            Часть маршрута «{route.parentRoute}» · {route.region}
          </p>
          <p className="mt-5 text-sm leading-7 text-brand-50/95 sm:text-base">{route.lead}</p>
          <p className="mt-4 text-sm leading-7 text-brand-100/90 sm:text-base">{route.body}</p>
          <div className="mt-6 flex flex-wrap items-center gap-3">
            <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-amber-100">
              {route.duration}
            </span>
            <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-brand-100">
              Выезд из Санкт-Петербурга
            </span>
          </div>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link href={route.ctaHref} className="btn-primary bg-white text-brand-900 hover:bg-amber-50">
              {route.ctaLabel}
            </Link>
            <Link
              href={route.secondaryHref}
              className="rounded-full border border-white/30 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-white/10"
            >
              {route.secondaryCta}
            </Link>
          </div>
        </div>

        <div className="px-6 py-8 sm:px-8 sm:py-10">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-brand-700">Программа</p>
          <ol className="mt-5 space-y-6">
            {route.days.map((day, index) => (
              <li key={day.title}>
                <p className="flex items-baseline gap-3 font-semibold text-stone-900">
                  <span className="font-display text-2xl text-brand-200">{String(index + 1).padStart(2, "0")}</span>
                  {day.title}
                </p>
                <ul className="mt-3 space-y-2 text-sm leading-6 text-stone-600">
                  {day.points.map((point) => (
                    <li key={point} className="flex gap-2">
                      <span className="mt-2 size-1.5 shrink-0 rounded-full bg-brand-700" />
                      <span>{point}</span>
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ol>
        </div>
      </div>
    </section>
  );
}
