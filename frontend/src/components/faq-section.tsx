"use client";

import { useState } from "react";
import { faqItems } from "@/lib/site-content";
import { SectionHeading } from "@/components/section-heading";
import type { FaqBlockContent } from "@/lib/api/cms";

type FaqSectionProps = {
  content?: FaqBlockContent;
};

export function FaqSection({ content }: FaqSectionProps = {}) {
  const [openIndex, setOpenIndex] = useState<number | null>(0);
  const eyebrow = content?.eyebrow ?? "Вопросы";
  const title = content?.title ?? "Частые вопросы";
  const description =
    content?.description ?? "Не нашли ответ? Позвоните нам — менеджер всё объяснит.";
  const items = content?.items?.length ? content.items : faqItems;

  return (
    <section id="faq" className="scroll-mt-24">
      {items.length > 0 ? (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "FAQPage",
              mainEntity: items.map((item) => ({
                "@type": "Question",
                name: item.question,
                acceptedAnswer: { "@type": "Answer", text: item.answer },
              })),
            }),
          }}
        />
      ) : null}
      <SectionHeading eyebrow={eyebrow} title={title} description={description} />

      <div className="mt-8 divide-y divide-stone-200 rounded-2xl border border-stone-200 bg-white">
        {items.map((item, index) => {
          const isOpen = openIndex === index;

          return (
            <div key={item.question}>
              <button
                type="button"
                onClick={() => setOpenIndex(isOpen ? null : index)}
                className="flex w-full items-center justify-between gap-4 px-5 py-4 text-left transition hover:bg-stone-50"
                aria-expanded={isOpen}
                aria-controls={`faq-answer-${index}`}
                id={`faq-question-${index}`}
              >
                <span className="font-medium text-stone-900">{item.question}</span>
                <span
                  className={`flex size-7 shrink-0 items-center justify-center rounded-full bg-stone-100 text-stone-600 transition ${isOpen ? "rotate-45" : ""}`}
                  aria-hidden
                >
                  +
                </span>
              </button>
              <div
                id={`faq-answer-${index}`}
                role="region"
                aria-labelledby={`faq-question-${index}`}
                className={`px-5 text-sm leading-7 text-stone-600 ${isOpen ? "pb-4" : "sr-only"}`}
              >
                {item.answer}
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}
