import Link from "next/link";
import type { Metadata } from "next";
import { SectionHeading } from "@/components/section-heading";
import { siteConfig } from "@/lib/site-config";

export const metadata: Metadata = {
  title: "Политика конфиденциальности",
  description: `Как ${siteConfig.name} обрабатывает персональные данные из заявок на сайте.`,
  alternates: { canonical: "/privacy" },
  openGraph: {
    title: "Политика конфиденциальности",
    description: `Как ${siteConfig.name} обрабатывает персональные данные из заявок на сайте.`,
    url: "/privacy",
  },
};

export default function PrivacyPage() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-10 sm:py-14">
      <div className="mb-6 flex flex-wrap gap-4 text-sm">
        <Link href="/" className="text-stone-500 transition hover:text-brand-800">
          ← На главную
        </Link>
        <Link href="/support" className="text-stone-500 transition hover:text-brand-800">
          Поддержка
        </Link>
        <Link href="/search" className="text-stone-500 transition hover:text-brand-800">
          Туры
        </Link>
      </div>

      <SectionHeading
        title="Политика обработки персональных данных"
        description="Краткая информация о том, как мы используем данные из заявок на сайте."
      />

      <div className="prose prose-stone mt-8 max-w-none text-sm leading-7 text-stone-700">
        <p>
          Оставляя заявку на сайте {siteConfig.fullName}, вы передаёте имя, телефон и, при желании, email и
          комментарий. Эти данные нужны только для связи по вашей поездке и организации участия.
        </p>
        <p>
          Мы не продаём персональные данные третьим лицам. Доступ к заявкам имеют только
          уполномоченные сотрудники службы.
        </p>
        <p>
          Данные хранятся в защищённой системе учёта заявок столько, сколько необходимо для
          обработки поездки и выполнения требований законодательства.
        </p>
        <p>
          По вопросам обработки данных вы можете связаться с нами по контактам, указанным в подвале
          сайта.
        </p>
        <p className="text-stone-500">
          Текст является базовой версией для MVP. Перед публичным запуском рекомендуется
          согласовать финальную редакцию с юристом.
        </p>
      </div>
    </div>
  );
}
