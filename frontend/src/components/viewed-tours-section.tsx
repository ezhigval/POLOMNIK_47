import Link from "next/link";
import { TourSchedule } from "@/components/tour-schedule";
import { SectionHeading } from "@/components/section-heading";
import { loadViewedTours } from "@/lib/viewed-tours-server";

export async function ViewedToursSection() {
  const tours = await loadViewedTours(5);
  if (tours.length === 0) {
    return null;
  }

  return (
    <section className="scroll-mt-24">
      <SectionHeading
        eyebrow="Недавно"
        title="Вы смотрели"
        description="Туры, которые вы недавно открывали на сайте."
      />
      <div className="mt-6">
        <TourSchedule tours={tours} />
      </div>
      <p className="mt-4 text-center">
        <Link href="/search" className="text-sm font-medium text-brand-800 hover:text-brand-900">
          Все туры в расписании →
        </Link>
      </p>
    </section>
  );
}
