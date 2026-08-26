import { SectionHeading } from "@/components/section-heading";
import { DioceseAffiliation } from "@/components/diocese-affiliation";
import { contactPhoneDisplay, contactEmail } from "@/lib/contact";
import { siteConfig } from "@/lib/site-config";
import type { AboutBlockContent } from "@/lib/api/cms";

const defaultHighlights = [
  "Групповые поездки с продуманной программой",
  "Сопровождение священника на маршруте",
  "Помощь менеджера на каждом этапе — от заявки до возвращения",
  `Выезд из ${siteConfig.departureCity} и других городов по согласованию`,
];

const defaultParagraphs = [
  "«Тихвинский путь» организует паломнические поездки по России: от классических направлений — Оптина пустынь, Дивеево, Валаам — до сезонных поездок в монастыри и святые места Северо-Запада.",
  "Наша задача — создать все условия для того, чтобы вы могли прикоснуться к православным святыням.",
  "Записаться просто: оставьте заявку на сайте, менеджер перезвонит, уточнит детали и подтвердит участие. Оплату на сайте вносить не нужно.",
];

type AboutSectionProps = {
  content?: AboutBlockContent;
};

export function AboutSection({ content }: AboutSectionProps = {}) {
  const eyebrow = content?.eyebrow ?? "О службе";
  const title = content?.title ?? siteConfig.fullName;
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
              <DioceseAffiliation
                className="mt-4"
                textClassName="font-medium text-stone-900"
                linkClassName="mt-1 inline-flex items-center text-brand-800 underline decoration-brand-200 underline-offset-2 hover:decoration-brand-800"
              />
            </div>
          ) : null}
        </div>
      </div>
    </section>
  );
}
