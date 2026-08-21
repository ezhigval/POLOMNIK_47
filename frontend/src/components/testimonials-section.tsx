import Link from "next/link";
import { SectionHeading } from "@/components/section-heading";
import { TestimonialCard } from "@/components/testimonial-card";
import { loadTestimonials } from "@/lib/testimonials";

export async function TestimonialsSection() {
  const testimonials = await loadTestimonials(3);

  return (
    <section id="reviews" className="scroll-mt-24">
      <SectionHeading
        eyebrow="Отзывы"
        title="Слова тех, кто уже с нами ездил"
        description="Реальные впечатления паломников после поездок."
      />

      <div className="mt-8 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {testimonials.map((item, index) => (
          <TestimonialCard key={`${item.client_name}-${index}`} item={item} />
        ))}
      </div>

      <div className="mt-10 flex justify-center">
        <Link href="/reviews" className="btn-secondary px-6 py-2.5 text-sm">
          Смотреть все
        </Link>
      </div>
    </section>
  );
}
