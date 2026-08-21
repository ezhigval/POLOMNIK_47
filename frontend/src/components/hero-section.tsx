import Link from "next/link";
import { OptimizedImage } from "@/components/optimized-image";
import { TripSearchConstructor } from "@/components/trip-search-constructor";
import { heroBackgroundImage } from "@/lib/tour-cover";
import { heroContent, trustStats } from "@/lib/site-content";
import type { HeroBlockContent } from "@/lib/api/cms";

type HeroSectionProps = {
  content?: HeroBlockContent;
};

export function HeroSection({ content }: HeroSectionProps = {}) {
  const eyebrow = content?.eyebrow || heroContent.eyebrow;
  const title = content?.title || heroContent.title;
  const subtitle = content?.subtitle || heroContent.subtitle;
  const stats = content?.stats?.length ? content.stats : trustStats;
  const primaryCta = content?.primaryCta || heroContent.primaryCta;
  const primaryHref = content?.primaryHref || "/search";
  const secondaryCta = content?.secondaryCta || heroContent.secondaryCta;
  const secondaryHref = content?.secondaryHref || "/#how-it-works";

  return (
    <section className="relative overflow-hidden rounded-3xl">
      <div className="absolute inset-0">
        <OptimizedImage
          src={heroBackgroundImage}
          alt=""
          fill
          priority
          sizes="100vw"
          className="object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-br from-brand-950/92 via-brand-900/85 to-brand-800/75" />
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_rgba(251,191,36,0.15),_transparent_50%)]" />
      </div>

      <div className="relative px-4 py-10 sm:px-8 sm:py-14 lg:py-16">
        <div className="max-w-3xl">
          <p className="mb-4 inline-flex items-center gap-2 rounded-full border border-white/20 bg-white/10 px-3 py-1 text-xs font-medium uppercase tracking-widest text-amber-100 backdrop-blur-sm">
            <span className="size-1.5 rounded-full bg-amber-300" />
            {eyebrow}
          </p>
          <h1 className="font-display text-4xl font-semibold leading-[1.08] tracking-tight text-white sm:text-5xl lg:text-6xl">
            {title}
          </h1>
          <p className="mt-4 max-w-2xl text-base leading-7 text-brand-50/95 sm:text-lg">
            {subtitle}
          </p>
          {primaryCta || secondaryCta ? (
            <div className="mt-6 flex flex-wrap gap-3">
              {primaryCta ? (
                <Link href={primaryHref} className="btn-primary bg-white text-brand-900 hover:bg-amber-50">
                  {primaryCta}
                </Link>
              ) : null}
              {secondaryCta ? (
                <Link
                  href={secondaryHref}
                  className="rounded-full border border-white/30 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-white/10"
                >
                  {secondaryCta}
                </Link>
              ) : null}
            </div>
          ) : null}
        </div>

        <div className="relative z-10 mt-8 lg:mt-10">
          <TripSearchConstructor className="shadow-2xl shadow-black/20" />
        </div>

        <dl className="mt-10 grid grid-cols-2 gap-3 sm:grid-cols-4 sm:gap-4">
          {stats.map((stat) => (
            <div
              key={stat.label}
              className="rounded-2xl border border-white/15 bg-white/10 px-4 py-3 backdrop-blur-sm"
            >
              <dt className="font-display text-2xl font-semibold text-white sm:text-3xl">{stat.value}</dt>
              <dd className="mt-0.5 text-xs text-brand-100 sm:text-sm">{stat.label}</dd>
            </div>
          ))}
        </dl>
      </div>
    </section>
  );
}
