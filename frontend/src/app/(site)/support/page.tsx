import Link from "next/link";
import type { Metadata } from "next";
import { DioceseAffiliation } from "@/components/diocese-affiliation";
import { FaqSection } from "@/components/faq-section";
import { SectionHeading } from "@/components/section-heading";
import { contactEmail, contactPhone, contactPhoneDisplay } from "@/lib/contact";

export const metadata: Metadata = {
  title: "Поддержка",
  description: "Помощь по бронированию, оплате, документам и поездке.",
  alternates: { canonical: "/support" },
};

const topics = [
  {
    title: "Бронирование",
    items: [
      "Как оставить заявку на тур",
      "Можно ли изменить количество человек",
      "Как узнать статус заявки",
    ],
  },
  {
    title: "Оплата",
    items: [
      "Когда и как оплачивать поездку",
      "Что входит в стоимость тура",
      "Возврат при отмене",
    ],
  },
  {
    title: "Поездка",
    items: [
      "Что взять с собой в паломничество",
      "Как добраться до места сбора",
      "Правила поведения в монастыре",
    ],
  },
  {
    title: "Личный кабинет",
    items: [
      "Регистрация и вход",
      "Где смотреть мои заявки",
      "Как связаться с менеджером",
    ],
  },
];

export default function SupportPage() {
  return (
    <div className="mx-auto max-w-6xl space-y-16 px-4 py-8 sm:py-12">
      <section className="rounded-3xl bg-gradient-to-br from-brand-950 via-brand-900 to-brand-800 px-6 py-10 text-white sm:px-10">
        <p className="text-sm font-medium uppercase tracking-widest text-brand-100">Поддержка</p>
        <h1 className="mt-3 font-display text-4xl font-semibold sm:text-5xl">Чем можем помочь?</h1>
        <p className="mt-4 max-w-2xl text-brand-50/90">
          Ответы на частые вопросы, контакты менеджера и помощь с подбором поездки.
        </p>
        <div className="mt-8 flex flex-wrap gap-3">
          <a href={`tel:${contactPhone}`} className="btn-primary bg-white text-brand-900 hover:bg-brand-50">
            {contactPhoneDisplay}
          </a>
          <a
            href={`mailto:${contactEmail}`}
            className="btn-secondary border-white/30 bg-white/10 text-white hover:bg-white/20"
          >
            {contactEmail}
          </a>
          <Link href="/support/chat" className="btn-secondary border-white/30 bg-white/10 text-white hover:bg-white/20">
            Открыть чат
          </Link>
        </div>
        <DioceseAffiliation
          className="mt-6"
          textClassName="text-sm text-brand-100/90"
          linkClassName="mt-2 inline-flex items-center text-sm text-white underline decoration-white/40 underline-offset-2 hover:decoration-white"
        />
      </section>

      <section>
        <SectionHeading
          eyebrow="Справочник"
          title="Популярные темы"
          description="Быстрые ответы по основным сценариям."
        />
        <div className="mt-8 grid gap-4 md:grid-cols-2">
          {topics.map((topic) => (
            <article
              key={topic.title}
              className="rounded-2xl border border-stone-200 bg-white p-5 shadow-sm"
            >
              <h2 className="font-semibold text-stone-900">{topic.title}</h2>
              <ul className="mt-3 space-y-2 text-sm text-stone-600">
                {topic.items.map((item) => (
                  <li key={item} className="flex gap-2">
                    <span className="text-brand-700">•</span>
                    {item}
                  </li>
                ))}
              </ul>
            </article>
          ))}
        </div>
      </section>

      <FaqSection />

      <section className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
          <h2 className="text-xl font-semibold text-stone-900">Нужна помощь с подбором?</h2>
          <p className="mt-2 text-sm leading-7 text-stone-600">
            Опишите желаемое направление, даты и состав группы — менеджер предложит варианты и
            перезвонит в рабочее время.
          </p>
          <Link href="/search" className="btn-primary mt-4 inline-flex">
            Открыть поиск туров
          </Link>
        </div>
        <div className="rounded-2xl border border-stone-200 bg-white p-6 shadow-sm">
          <h2 className="text-xl font-semibold text-stone-900">Уже есть аккаунт?</h2>
          <p className="mt-2 text-sm leading-7 text-stone-600">
            Войдите, чтобы видеть статус заявок и историю поездок в одном месте.
          </p>
          <div className="mt-4 flex flex-wrap gap-3">
            <Link href="/account/login" className="btn-primary">
              Войти
            </Link>
            <Link href="/account/register" className="btn-secondary">
              Регистрация
            </Link>
            <Link href="/support/chat" className="btn-secondary">
              Чат поддержки
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}
