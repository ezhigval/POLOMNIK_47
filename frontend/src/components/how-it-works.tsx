import Link from "next/link";
import { SectionHeading } from "@/components/section-heading";
import type { HowItWorksBlockContent } from "@/lib/api/cms";

const defaultSteps = [
  {
    title: "Выберите тур",
    description: "Изучите даты, программу и стоимость. Фильтры помогут найти подходящее направление.",
  },
  {
    title: "Оставьте заявку",
    description: "Укажите имя, телефон и количество участников — без регистрации и оплаты на сайте.",
  },
  {
    title: "Подтверждение",
    description: "Менеджер перезвонит, ответит на вопросы и подтвердит ваше участие в группе.",
  },
];

type HowItWorksSectionProps = {
  content?: HowItWorksBlockContent;
};

export function HowItWorksSection({ content }: HowItWorksSectionProps = {}) {
  const eyebrow = content?.eyebrow ?? "Просто";
  const title = content?.title ?? "Как записаться";
  const description =
    content?.description ??
    "Три шага — и вы в списке участников. Никаких личных кабинетов и сложных форм.";
  const steps = content?.steps?.length ? content.steps : defaultSteps;
  const ctaLabel = content?.ctaLabel ?? "Записаться";
  const ctaHref = content?.ctaHref ?? "/search";

  return (
    <section id="how-it-works" className="scroll-mt-24">
      <SectionHeading eyebrow={eyebrow} title={title} description={description} />

      <ol className="mt-10 grid gap-5 sm:grid-cols-3">
        {steps.map((step, index) => (
          <li
            key={step.title}
            className="relative rounded-2xl border border-stone-200 bg-white p-6 shadow-sm"
          >
            {index < steps.length - 1 ? (
              <span
                className="absolute -right-3 top-1/2 hidden size-6 -translate-y-1/2 items-center justify-center rounded-full bg-brand-100 text-brand-800 sm:flex"
                aria-hidden
              >
                →
              </span>
            ) : null}
            <span className="font-display text-3xl font-semibold text-brand-200">
              {String(index + 1).padStart(2, "0")}
            </span>
            <h3 className="mt-3 mb-2 text-lg font-semibold text-stone-900">{step.title}</h3>
            <p className="text-sm leading-6 text-stone-600">{step.description}</p>
          </li>
        ))}
      </ol>

      <div className="mt-10 flex justify-center">
        <Link href={ctaHref} className="btn-primary px-8 py-3 text-base">
          {ctaLabel}
        </Link>
      </div>
    </section>
  );
}
