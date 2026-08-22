import { Suspense } from "react";
import { AboutSection } from "@/components/about-section";
import { CtaSection } from "@/components/cta-section";
import { FaqSection } from "@/components/faq-section";
import { FeaturedRouteSection } from "@/components/featured-route-section";
import { HeroSection } from "@/components/hero-section";
import { HowItWorksSection } from "@/components/how-it-works";
import { PopularDestinations, PopularDestinationsSkeleton } from "@/components/popular-destinations";
import { SectionHeading } from "@/components/section-heading";
import { TestimonialsSection } from "@/components/testimonials-section";
import { WhyUsSection } from "@/components/why-us-section";
import type {
  AboutBlockContent,
  CmsBlock,
  CtaBlockContent,
  FaqBlockContent,
  FeaturedRouteBlockContent,
  HeroBlockContent,
  HowItWorksBlockContent,
  RichTextBlockContent,
  WhyUsBlockContent,
} from "@/lib/api/cms";

type CmsPageRendererProps = {
  blocks: CmsBlock[];
};

function asContent<T>(content: Record<string, unknown>): T {
  return content as T;
}

function firstHeroStats(blocks: CmsBlock[]) {
  for (const block of blocks) {
    if (block.type !== "hero") continue;
    const stats = asContent<HeroBlockContent>(block.content).stats;
    if (stats?.length) return stats;
  }
  return undefined;
}

export function CmsPageRenderer({ blocks }: CmsPageRendererProps) {
  const ordered = [...blocks].sort((a, b) => a.sort_order - b.sort_order);
  const heroBlocks = ordered.filter((block) => block.type === "hero");
  const bodyBlocks = ordered.filter((block) => block.type !== "hero");
  const hasFeaturedRoute = ordered.some((block) => block.type === "featured_route" && block.is_visible);
  const heroStats = firstHeroStats(heroBlocks);

  return (
    <>
      {heroBlocks.map((block) => (
        <div key={block.id} className="mx-auto max-w-6xl px-4 py-6 sm:py-8">
          <HeroSection content={asContent<HeroBlockContent>(block.content)} />
        </div>
      ))}

      {!hasFeaturedRoute ? (
        <div className="mx-auto max-w-6xl px-4 pb-8 sm:pb-12">
          <FeaturedRouteSection />
        </div>
      ) : null}

      {bodyBlocks.length > 0 ? (
        <div className="mx-auto max-w-6xl space-y-20 px-4 py-8 sm:space-y-24 sm:py-12">
          {bodyBlocks.map((block) => (
            <CmsBlockView key={block.id} block={block} heroStats={heroStats} />
          ))}
        </div>
      ) : null}
    </>
  );
}

function CmsBlockView({
  block,
  heroStats,
}: {
  block: CmsBlock;
  heroStats?: HeroBlockContent["stats"];
}) {
  switch (block.type) {
    case "about":
      return <AboutSection content={asContent<AboutBlockContent>(block.content)} />;
    case "why_us":
      return (
        <WhyUsSection content={asContent<WhyUsBlockContent>(block.content)} fallbackStats={heroStats} />
      );
    case "how_it_works":
      return <HowItWorksSection content={asContent<HowItWorksBlockContent>(block.content)} />;
    case "faq":
      return <FaqSection content={asContent<FaqBlockContent>(block.content)} />;
    case "cta":
      return <CtaSection content={asContent<CtaBlockContent>(block.content)} />;
    case "rich_text": {
      const content = asContent<RichTextBlockContent>(block.content);
      return (
        <section className="scroll-mt-24">
          <SectionHeading
            eyebrow={content.eyebrow}
            title={content.title || "Текст"}
          />
          {content.body ? (
            <div className="mt-6 whitespace-pre-wrap text-sm leading-7 text-stone-700">
              {content.body}
            </div>
          ) : null}
        </section>
      );
    }
    case "featured_route":
      return <FeaturedRouteSection content={asContent<FeaturedRouteBlockContent>(block.content)} />;
    case "popular_destinations":
      return (
        <Suspense fallback={<PopularDestinationsSkeleton />}>
          <PopularDestinations />
        </Suspense>
      );
    case "testimonials":
      return (
        <Suspense fallback={<ReviewsSkeleton />}>
          <TestimonialsSection />
        </Suspense>
      );
    case "hero":
      return <HeroSection content={asContent<HeroBlockContent>(block.content)} />;
    default:
      return null;
  }
}

function ReviewsSkeleton() {
  return (
    <div className="grid gap-4 md:grid-cols-3" aria-busy="true">
      {Array.from({ length: 3 }).map((_, i) => (
        <div key={i} className="h-48 animate-pulse rounded-2xl bg-stone-200" />
      ))}
    </div>
  );
}
