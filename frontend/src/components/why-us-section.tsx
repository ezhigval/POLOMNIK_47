import type { ReactNode } from "react";
import { trustStats, whyUsItems } from "@/lib/site-content";
import { SectionHeading } from "@/components/section-heading";
import type { WhyUsBlockContent } from "@/lib/api/cms";

const icons: Record<string, ReactNode> = {
  route: (
    <svg viewBox="0 0 24 24" fill="none" className="size-6" stroke="currentColor" strokeWidth="1.5">
      <path strokeLinecap="round" strokeLinejoin="round" d="M9 6.75V15m0 0H6.75m2.25 0h2.25M9 15v3.75M15 9.75V18m0 0h-2.25M15 18h2.25M15 9.75V6m0 3.75h2.25M15 6h-2.25" />
    </svg>
  ),
  cross: (
    <svg viewBox="0 0 24 24" fill="none" className="size-6" stroke="currentColor" strokeWidth="1.5">
      <path strokeLinecap="round" d="M12 4v16M8 8h8" />
    </svg>
  ),
  shield: (
    <svg viewBox="0 0 24 24" fill="none" className="size-6" stroke="currentColor" strokeWidth="1.5">
      <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z" />
    </svg>
  ),
  wallet: (
    <svg viewBox="0 0 24 24" fill="none" className="size-6" stroke="currentColor" strokeWidth="1.5">
      <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a2.25 2.25 0 0 0-2.25-2.25H15a3 3 0 1 1-6 0H5.25A2.25 2.25 0 0 0 3 12m18 0v6a2.25 2.25 0 0 1-2.25 2.25H5.25A2.25 2.25 0 0 1 3 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 0 0-2.25-2.25H5.25A2.25 2.25 0 0 0 3 9m18 0V6a2.25 2.25 0 0 0-2.25-2.25H5.25A2.25 2.25 0 0 0 3 6v3" />
    </svg>
  ),
};

type WhyUsStat = { value: string; label: string };

type WhyUsSectionProps = {
  content?: WhyUsBlockContent;
  fallbackStats?: WhyUsStat[];
};

function statsWithoutOverlap(stats: WhyUsStat[], items: { title: string; description: string }[]) {
  return stats.filter((stat) => {
    const value = stat.value.trim().toLowerCase();
    const label = stat.label.trim().toLowerCase();
    return !items.some((item) => {
      const text = `${item.title} ${item.description}`.toLowerCase();
      return Boolean(value && label && text.includes(value) && text.includes(label));
    });
  });
}

export function WhyUsSection({ content, fallbackStats }: WhyUsSectionProps = {}) {
  const eyebrow = content?.eyebrow ?? "Почему мы";
  const title = content?.title ?? "Паломничество без лишних забот";
  const description =
    content?.description?.trim() ||
    "Мы не просто везём вас в монастырь — мы создаём возможность для тишины, молитвы и встречи с Богом.";
  const items = content?.items?.length ? content.items : whyUsItems;
  const stats = statsWithoutOverlap(
    content?.stats?.length ? content.stats : fallbackStats?.length ? fallbackStats : trustStats,
    items,
  );

  return (
    <section id="why-us" className="scroll-mt-24">
      <SectionHeading eyebrow={eyebrow} title={title} description={description} />

      {stats.length ? (
        <dl className="mt-10 grid grid-cols-2 gap-5 sm:grid-cols-4">
          {stats.map((stat) => (
            <div
              key={stat.label}
              className="rounded-2xl border border-stone-200/80 bg-white px-4 py-5 shadow-sm transition hover:border-brand-200 hover:shadow-md"
            >
              <dt className="font-display text-2xl font-semibold text-stone-900 sm:text-3xl">{stat.value}</dt>
              <dd className="mt-1 text-xs text-stone-600 sm:text-sm">{stat.label}</dd>
            </div>
          ))}
        </dl>
      ) : null}

      <ul className="mt-10 grid gap-5 sm:grid-cols-2">
        {items.map((item) => (
          <li
            key={item.title}
            className="group rounded-2xl border border-stone-200/80 bg-white p-6 shadow-sm transition hover:border-brand-200 hover:shadow-md"
          >
            <div className="mb-4 inline-flex size-11 items-center justify-center rounded-xl bg-brand-50 text-brand-800 transition group-hover:bg-brand-100">
              {icons[item.icon] ?? icons.route}
            </div>
            <h3 className="mb-2 text-lg font-semibold text-stone-900">{item.title}</h3>
            <p className="text-sm leading-6 text-stone-600">{item.description}</p>
          </li>
        ))}
      </ul>
    </section>
  );
}
