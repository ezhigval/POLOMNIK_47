import Link from "next/link";
import { ctaBanner } from "@/lib/site-content";
import type { CtaBlockContent } from "@/lib/api/cms";

type CtaSectionProps = {
  content?: CtaBlockContent;
};

export function CtaSection({ content }: CtaSectionProps = {}) {
  const title = content?.title ?? ctaBanner.title;
  const subtitle = content?.subtitle ?? ctaBanner.subtitle;
  const button = content?.button ?? ctaBanner.button;
  const href = content?.href ?? "/search";

  return (
    <section className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-900 via-brand-800 to-brand-700 px-6 py-12 text-center sm:px-10 sm:py-14">
      <div className="pointer-events-none absolute -right-20 -top-20 size-64 rounded-full bg-amber-300/20 blur-3xl" />
      <div className="pointer-events-none absolute -bottom-16 -left-16 size-56 rounded-full bg-white/10 blur-3xl" />

      <div className="relative mx-auto max-w-xl">
        <h2 className="font-display text-3xl font-semibold tracking-tight text-white sm:text-4xl">
          {title}
        </h2>
        <p className="mt-4 text-base leading-7 text-brand-100">{subtitle}</p>
        <Link
          href={href}
          className="btn-primary mt-8 bg-white text-brand-900 hover:bg-amber-50"
        >
          {button}
        </Link>
      </div>
    </section>
  );
}
