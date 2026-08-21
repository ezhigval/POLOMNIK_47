import { SectionHeading } from "@/components/section-heading";
import { contactPhoneDisplay, contactEmail } from "@/lib/contact";
import { siteConfig } from "@/lib/site-config";
import type { AboutBlockContent } from "@/lib/api/cms";

const defaultHighlights = [
  "Групповые поездки с продуманной программой",
  "Сопровождение духовника на маршруте",
  "Помощь менеджера на каждом этапе — от заявки до возвращения",
  `Выезд из ${siteConfig.departureCity} и других городов по согласованию`,
];

const defaultParagraphs = [
  `${siteConfig.name} организует паломнические туры по России: от классических направлений — Оптина, Дивеево, Валаам — до сезонных поездок в монастыри и святые места Северо-Запада.`,
  "Наша задача — создать спокойное пространство для паломничества: комфортный транспорт, проверенное размещение, понятная программа и человек на связи, если возникнут вопросы.",
  "Запись через сайт простая: вы оставляете заявку, менеджер перезванивает, уточняет детали и подтверждает участие. Без личного кабинета и онлайн-оплаты на первом этапе — всё через живой контакт.",
];

type AboutSectionProps = {
  content?: AboutBlockContent;
};

export function AboutSection({ content }: AboutSectionProps = {}) {
  const eyebrow = content?.eyebrow ?? "О службе";
  const title =
    content?.title ?? `${siteConfig.name} — паломничество без лишней суеты`;
  const paragraphs = content?.paragraphs?.filter(Boolean).length
    ? content.paragraphs.filter(Boolean)
    : defaultParagraphs;
  const highlights = content?.highlights?.length ? content.highlights : defaultHighlights;
  const showContacts = content?.showContacts ?? true;

  return (
    <section id="about" className="scroll-mt-24">
      <SectionHeading
        eyebrow={eyebrow}
        title={title}
        description={
          content
            ? undefined
            : "Мы помогаем людям доехать до святынь, сосредоточиться на молитве и не думать о логистике."
        }
      />

      <div className="mt-8 grid gap-8 lg:grid-cols-2">
        <div className="space-y-4 text-sm leading-7 text-stone-700">
          {paragraphs.map((paragraph, index) => (
            <p key={`${index}-${paragraph.slice(0, 24)}`}>{paragraph}</p>
          ))}
        </div>

        <div className="rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
          <h3 className="mb-4 font-semibold text-stone-900">Почему нам доверяют</h3>
          <ul className="space-y-3">
            {highlights.map((item) => (
              <li key={item} className="flex gap-3 text-sm text-stone-700">
                <span className="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-brand-50 text-xs text-brand-800">
                  ✓
                </span>
                {item}
              </li>
            ))}
          </ul>
          {showContacts ? (
            <div className="mt-6 border-t border-stone-100 pt-4 text-sm text-stone-600">
              <p className="font-medium text-stone-900">Связаться с нами</p>
              <p className="mt-1">{contactPhoneDisplay}</p>
              <p>{contactEmail}</p>
            </div>
          ) : null}
        </div>
      </div>
    </section>
  );
}
