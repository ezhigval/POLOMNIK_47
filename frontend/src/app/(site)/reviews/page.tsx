import type { Metadata } from "next";
import Link from "next/link";
import { PageIntro } from "@/components/page-intro";
import { PublicReviewForm } from "@/components/public-review-form";
import { TestimonialCard } from "@/components/testimonial-card";
import { getTours } from "@/lib/api/tours";
import { siteConfig } from "@/lib/site-config";
import { loadTestimonials } from "@/lib/testimonials";
import { buildPublicPageMetadata } from "@/lib/seo-metadata";

const reviewsDescription = `Отзывы паломников о поездках «${siteConfig.name}».`;

export const metadata: Metadata = buildPublicPageMetadata({
  title: "Отзывы",
  description: reviewsDescription,
  canonical: "/reviews",
});

export default async function ReviewsPage() {
  const testimonials = await loadTestimonials(24);
  let tours: { id: string; title: string }[] = [];
  try {
    const list = await getTours();
    tours = list.data.map((tour) => ({ id: tour.id, title: tour.title }));
  } catch {
    tours = [];
  }

  return (
    <div className="mx-auto max-w-6xl space-y-8 px-4 py-8 sm:py-10">
      <PageIntro
        title="Что говорят паломники"
        description="Впечатления после поездок: организация, сопровождение и атмосфера маршрута."
      />

      {testimonials.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-stone-300 bg-white p-8 text-center text-stone-500">
          Отзывов пока нет.{" "}
          <Link href="/search" className="font-medium text-brand-800 underline-offset-2 hover:underline">
            Выберите тур
          </Link>
          .
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {testimonials.map((item, index) => (
            <TestimonialCard key={`${item.client_name}-${index}`} item={item} />
          ))}
        </div>
      )}

      <div className="mx-auto max-w-xl">
        <PublicReviewForm tours={tours} />
      </div>

      <div className="flex justify-center">
        <Link href="/search" className="btn-primary px-8 py-3">
          Выбрать тур
        </Link>
      </div>
    </div>
  );
}
