import { HeroBackground } from "@/components/hero-background";
import { TripSearchConstructor } from "@/components/trip-search-constructor";
import { heroContent } from "@/lib/site-content";
import type { HeroBlockContent } from "@/lib/api/cms";

type HeroSectionProps = {
  content?: HeroBlockContent;
};

export function HeroSection({ content }: HeroSectionProps = {}) {
  const eyebrow = content?.eyebrow || heroContent.eyebrow;
  const title = content?.title || heroContent.title;
  const subtitle = content?.subtitle || heroContent.subtitle;

  return (
    <section className="relative overflow-hidden rounded-3xl">
      <HeroBackground />

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
        </div>

        <div className="relative z-10 mt-8 lg:mt-10">
          <TripSearchConstructor className="shadow-2xl shadow-black/20" />
        </div>
      </div>
    </section>
  );
}
