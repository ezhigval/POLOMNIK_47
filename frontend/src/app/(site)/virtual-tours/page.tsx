import type { Metadata } from "next";
import Link from "next/link";
import { Breadcrumbs } from "@/components/breadcrumbs";
import { PageIntro } from "@/components/page-intro";
import { SectionHeading } from "@/components/section-heading";
import { YandexPanoramaEmbed } from "@/components/virtual-tours/yandex-panorama-embed";
import { buildPublicPageMetadata } from "@/lib/seo-metadata";

export const metadata: Metadata = buildPublicPageMetadata({
  title: "3D-экскурсии",
  description: "Виртуальные панорамные экскурсии по святым местам и историческим локациям.",
  canonical: "/virtual-tours",
});

const panoramaItems = [
  {
    title: "Тихвинский монастырь",
    src: "https://yandex.ru/map-widget/v1/?ll=33.519991%2C59.653150&panorama%5Bdirection%5D=89.803885%2C12.604611&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.519991%2C59.653150&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Троицкий собор (Свято-Троицкий Александра Свирского мужской монастырь)",
    src: "https://yandex.ru/map-widget/v1/?ll=33.316935%2C60.779782&panorama%5Bdirection%5D=132.078479%2C7.885032&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.316437%2C60.779772&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Собор Преображения Господня (Свято-Троицкий Александра Свирского мужской монастырь)",
    src: "https://yandex.ru/map-widget/v1/?ll=33.309178%2C60.779273&panorama%5Bdirection%5D=37.484655%2C15.236205&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.309071%2C60.779156&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Мощи святого Александра Свирского (Свято-Троицкий Александра Свирского мужской монастырь)",
    src: "https://yandex.ru/map-widget/v1/?ll=33.309499%2C60.779429&panorama%5Bdirection%5D=171.379443%2C-3.596833&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.309482%2C60.779445&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Место явления Святой Троицы Александру Свирскому (Свято-Троицкий Александра Свирского мужской монастырь)",
    src: "https://yandex.ru/map-widget/v1/?ll=33.310736%2C60.779438&panorama%5Bdirection%5D=89.194018%2C-3.442152&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.310736%2C60.779438&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Часовня над источником Александра Свирского (Свято-Троицкий Александра Свирского мужской монастырь)",
    src: "https://yandex.ru/map-widget/v1/?ll=33.311548%2C60.779158&panorama%5Bdirection%5D=122.652909%2C5.061392&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.311369%2C60.779237&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Преображенская часть монастыря (Свято-Троицкий Александра Свирского мужской монастырь)",
    src: "https://yandex.ru/map-widget/v1/?ll=33.310450%2C60.779205&panorama%5Bdirection%5D=271.607362%2C8.530656&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=33.310450%2C60.779205&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
  {
    title: "Церковь Рождества Христова (пер. Челиева, 10)",
    src: "https://yandex.ru/map-widget/v1/?ll=30.467815%2C59.887786&panorama%5Bdirection%5D=54.882606%2C13.694496&panorama%5Bfull%5D=true&panorama%5Bpoint%5D=30.467438%2C59.887698&panorama%5Bspan%5D=82.776742%2C60.000000&z=10",
  },
];

export default function VirtualToursPage() {
  return (
    <div className="mx-auto max-w-6xl space-y-10 px-4 py-8 sm:py-10">
      <Breadcrumbs
        items={[
          { name: "Главная", href: "/" },
          { name: "3D-экскурсии" },
        ]}
      />

      <PageIntro
        showBack={false}
        title="3D-экскурсии"
        description="Виртуальные панорамы по святым местам и локациям маршрутов — чтобы заранее познакомиться с атмосферой поездки."
      />

      <section className="rounded-[2rem] border border-stone-200 bg-white p-4 shadow-sm sm:p-6 lg:p-8">
        <div className="mb-6 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-brand-700">Панорамы</p>
            <h2 className="mt-2 font-display text-3xl font-semibold text-stone-900 sm:text-4xl">
              Святые места и исторические точки
            </h2>
          </div>
          <Link href="/search" className="btn-primary inline-flex">
            Посмотреть поездки
          </Link>
        </div>

        <div className="space-y-8">
          {panoramaItems.map((item) => (
            <article key={item.title} className="rounded-[1.5rem] border border-stone-200 bg-stone-50 p-3 sm:p-4">
              <h3 className="mb-4 text-lg font-semibold text-stone-900 sm:text-xl">{item.title}</h3>
              <YandexPanoramaEmbed src={item.src} title={item.title} />
            </article>
          ))}
        </div>
      </section>

      <section className="grid gap-6 lg:grid-cols-[1.4fr_0.6fr]">
        <div className="rounded-[2rem] border border-stone-200 bg-white p-6 shadow-sm">
          <SectionHeading
            eyebrow="Описание"
            title="Как использовать"
            description="Страница собрана в расширяемом формате: для каждой локации достаточно добавить одну запись в массив и она автоматически появляется на странице."
          />
          <div className="mt-6 space-y-4 text-sm leading-7 text-stone-600">
            <p>
              Виртуальные экскурсии помогают заранее познакомиться с маршрутом, понять, какие места
              ждут паломников, и настроиться на поездку духовно и эмоционально.
            </p>
            <p>
              При необходимости можно быстро добавить новые точки, изменить заголовки или заменить
              embed-коды без правки разметки страницы.
            </p>
          </div>
        </div>

        <aside className="rounded-[2rem] border border-brand-100 bg-brand-50 p-6 shadow-sm">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-brand-700">Готово к расширению</p>
          <h3 className="mt-3 font-display text-2xl font-semibold text-brand-950">Следующие точки</h3>
          <ul className="mt-4 space-y-3 text-sm text-stone-700">
            <li className="flex gap-2"><span className="text-brand-700">•</span> Тихвин</li>
            <li className="flex gap-2"><span className="text-brand-700">•</span> Александро-Свирский монастырь</li>
            <li className="flex gap-2"><span className="text-brand-700">•</span> Исторические храмы Санкт-Петербурга</li>
          </ul>
        </aside>
      </section>
    </div>
  );
}
