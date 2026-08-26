import { Suspense } from "react";
import { AboutSection } from "@/components/about-section";
import { CmsPageRenderer } from "@/components/cms/cms-page-renderer";
import { CtaSection } from "@/components/cta-section";
import { FaqSection } from "@/components/faq-section";
import { FeaturedRouteSection } from "@/components/featured-route-section";
import { HeroSection } from "@/components/hero-section";
import { HomeNewsSection, HomeNewsSkeleton } from "@/components/home-news-section";
import { HowItWorksSection } from "@/components/how-it-works";
import { PopularDestinations, PopularDestinationsSkeleton } from "@/components/popular-destinations";
import { ViewedToursSection } from "@/components/viewed-tours-section";
import { TestimonialsSection } from "@/components/testimonials-section";
import { WhyUsSection } from "@/components/why-us-section";
import { getPublishedPage } from "@/lib/api/cms";
import { siteConfig } from "@/lib/site-config";
import type { Metadata } from "next";

export async function generateMetadata(): Promise<Metadata> {
  const cmsPage = await getPublishedPage("home").catch(() => null);
  const title = cmsPage?.meta_title?.trim() || undefined;
  const description =
    cmsPage?.meta_description?.trim() || siteConfig.description;

  return {
    ...(title ? { title: { absolute: title } } : {}),
    description,
    alternates: { canonical: "/" },
    openGraph: {
      title: title || `${siteConfig.name} — ${siteConfig.tagline}`,
      description,
      url: "/",
      type: "website",
    },
    twitter: {
      card: "summary_large_image",
      title: title || `${siteConfig.name} — ${siteConfig.tagline}`,
      description,
    },
  };
}

export default async function HomePage() {
  const cmsPage = await getPublishedPage("home").catch(() => null);

  if (cmsPage && cmsPage.blocks.length > 0) {
    return (
      <CmsPageRenderer
        blocks={cmsPage.blocks}
        afterHero={
          <Suspense fallback={<HomeNewsSkeleton />}>
            <HomeNewsSection />
          </Suspense>
        }
      />
    );
  }

  return (
    <>
      <div className="mx-auto max-w-6xl px-4 py-6 sm:py-8">
        <HeroSection />
      </div>

      <div className="mx-auto max-w-6xl space-y-20 px-4 py-8 sm:space-y-24 sm:py-12">
        <FeaturedRouteSection />
        <Suspense fallback={<HomeNewsSkeleton />}>
          <HomeNewsSection />
        </Suspense>
        <Suspense fallback={<PopularDestinationsSkeleton />}>
          <PopularDestinations />
        </Suspense>
        <ViewedToursSection />
        <AboutSection />
        <WhyUsSection />
        <HowItWorksSection />

        <Suspense fallback={<ReviewsSkeleton />}>
          <TestimonialsSection />
        </Suspense>

        <FaqSection />
        <CtaSection />
      </div>
    </>
  );
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
