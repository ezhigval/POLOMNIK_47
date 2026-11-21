import Link from "next/link";
import { russianPhone } from "@/lib/ru-contact";
import { OptimizedImage } from "@/components/optimized-image";
import { portraitUrl } from "@/lib/portrait-url";
import { Suspense, type ReactNode } from "react";
import Image from "next/image";
import { HomeNewsSection, HomeNewsSkeleton } from "@/components/home-news-section";

import { SectionHeading } from "@/components/section-heading";
import type {
  CmsBlock,
  HolyPlacesBlockContent,
  InformationHeroBlockContent,
  LaborTripsBlockContent,
  LegalDonationBlockContent,
  ServiceProfileBlockContent,
  ServiceTeamBlockContent,
} from "@/lib/api/cms";

export const informationBlockTypes = new Set([
  "information_hero",
  "holy_places",
  "service_profile",
  "service_team",
  "legal_donation",
  "labor_trips",
]);

function asContent<T>(content: Record<string, unknown>): T {
  return content as T;
}

function safeHref(value?: string): string | null {
  const href = value?.trim();
  if (!href || href.startsWith("//")) return null;
  if (href.startsWith("/") || href.startsWith("#")) return href;
  try {
    const url = new URL(href);
    return ["http:", "https:", "mailto:", "tel:"].includes(url.protocol) ? href : null;
  } catch {
    return null;
  }
}

function ContentLink({
  href,
  children,
  className,
}: {
  href?: string;
  children: ReactNode;
  className?: string;
}) {
  const target = safeHref(href);
  if (!target) return null;
  const external = target.startsWith("http:") || target.startsWith("https:");
  if (target.startsWith("/")) {
    return (
      <Link href={target} className={className}>
        {children}
      </Link>
    );
  }
  return (
    <a
      href={target}
      className={className}
      target={external ? "_blank" : undefined}
      rel={external ? "noopener noreferrer" : undefined}
    >
      {children}
    </a>
  );
}

export function InformationSpace({ blocks, fallback }: { blocks: CmsBlock[]; fallback: boolean }) {
  const ordered = blocks
    .filter((block) => block.is_visible && informationBlockTypes.has(block.type))
    .sort((a, b) => a.sort_order - b.sort_order);
  const hasHero = ordered.some((block) => block.type === "information_hero");

  return (
    <>
      {ordered.map((block) => {
        if (block.type === "information_hero") {
          return (
            <div key={block.id} className="mx-auto max-w-6xl px-4 py-6 sm:py-8">
              <InformationHero content={asContent<InformationHeroBlockContent>(block.content)} />
            </div>
          );
        }

        const contentClass = "mx-auto max-w-6xl px-4 py-8 sm:py-10";
        switch (block.type) {
          case "holy_places":
            return (
              <div key={block.id} className={contentClass}>
                <HolyPlaces content={asContent<HolyPlacesBlockContent>(block.content)} />
              </div>
            );
          case "service_profile":
            return (
              <div key={block.id} className={contentClass}>
                <ServiceProfile content={asContent<ServiceProfileBlockContent>(block.content)} />
              </div>
            );
          case "service_team":
            return (
              <div key={block.id} className={contentClass}>
                <ServiceTeam content={asContent<ServiceTeamBlockContent>(block.content)} />
              </div>
            );
          case "legal_donation":
            return (
              <div key={block.id} className={contentClass}>
                <LegalDonation content={asContent<LegalDonationBlockContent>(block.content)} />
              </div>
            );
          case "labor_trips":
            return (
              <div key={block.id} className={contentClass}>
                <LaborTrips content={asContent<LaborTripsBlockContent>(block.content)} />
              </div>
            );
          default:
            return null;
        }
      })}

      {fallback && !hasHero ? (
        <div className="mx-auto max-w-6xl px-4 py-6 sm:py-8">
          <InformationHero
            content={{
              eyebrow: "Паломническая служба Тихвинской епархии",
              title: "Дорога к святыням начинается с доброй встречи",
              subtitle: "Информация о службе, монастырях и паломнических поездках.",
              primaryCta: "Смотреть поездки",
              primaryHref: "/search",
              secondaryCta: "О службе",
              secondaryHref: "#service",
            }}
          />
        </div>
      ) : null}

      <div className="mx-auto max-w-6xl space-y-14 px-4 py-8 sm:space-y-16 sm:py-12">
        <Suspense fallback={<HomeNewsSkeleton />}>
          <HomeNewsSection />
        </Suspense>

      </div>
    </>
  );
}

function InformationHero({ content }: { content: InformationHeroBlockContent }) {
  return (
    <section className="relative min-h-[26rem] overflow-hidden rounded-3xl sm:min-h-[30rem]">
      <Image src="/images/hero/tikhvin-monastery.webp" alt="" fill priority sizes="(min-width: 1152px) 1120px, 100vw" className="object-cover" />
      <div className="absolute inset-0 bg-gradient-to-r from-black/80 via-black/50 to-black/10" />
      <div className="relative flex min-h-[26rem] items-end px-5 py-8 sm:min-h-[30rem] sm:px-10 sm:py-12 lg:px-14">
        <div className="max-w-3xl">
          {content.eyebrow ? (
            <p className="mb-4 inline-flex rounded-full border border-white/20 bg-white/10 px-3 py-1 text-xs font-medium uppercase tracking-widest text-amber-100 backdrop-blur-sm">
              {content.eyebrow}
            </p>
          ) : null}
          <h1 className="font-display text-4xl font-semibold leading-[1.08] tracking-tight text-white sm:text-5xl lg:text-6xl">
            {content.title}
          </h1>
          {content.subtitle ? (
            <p className="mt-5 max-w-2xl text-base leading-7 text-brand-50/95 sm:text-lg">
              {content.subtitle}
            </p>
          ) : null}
          <div className="mt-7 flex flex-wrap gap-3">
            {content.primaryCta ? (
              <ContentLink
                href={content.primaryHref}
                className="btn-primary bg-amber-300 text-brand-950 hover:bg-amber-200"
              >
                {content.primaryCta}
              </ContentLink>
            ) : null}
            {content.secondaryCta ? (
              <ContentLink
                href={content.secondaryHref}
                className="inline-flex items-center justify-center rounded-full border border-white/50 bg-white/10 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-white/20"
              >
                {content.secondaryCta}
              </ContentLink>
            ) : null}
          </div>
        </div>
      </div>
    </section>
  );
}

function HolyPlaces({ content }: { content: HolyPlacesBlockContent }) {
  const places = (content.items ?? []).filter((place) => place.name?.trim());
  if (places.length === 0) return null;

  return (
    <section className="scroll-mt-24 space-y-6" id="holy-places">
      <SectionHeading eyebrow={content.eyebrow} title={content.title || "Монастыри и святыни"} />
      {content.description ? <p className="max-w-3xl text-sm leading-7 text-stone-600">{content.description}</p> : null}
      <div className="grid gap-5 md:grid-cols-2">
        {places.map((place, index) => {
          const hasCoordinates = Boolean(place.latitude?.trim() && place.longitude?.trim());
          const latitude = Number(place.latitude);
          const longitude = Number(place.longitude);
          const coordinates = hasCoordinates && Math.abs(latitude) <= 90 && Math.abs(longitude) <= 180
            ? `${latitude.toFixed(5)}, ${longitude.toFixed(5)}`
            : "";
          const mapUrl = safeHref(place.mapUrl) || (coordinates
            ? `https://yandex.ru/maps/?pt=${longitude},${latitude}&z=15&l=map`
            : null);
          return (
            <article key={`${place.name}-${index}`} className="flex flex-col rounded-3xl border border-stone-200 bg-white p-6 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h3 className="font-display text-xl font-semibold text-stone-900">{place.name}</h3>
                  {place.region ? <p className="mt-1 text-sm font-medium text-brand-800">{place.region}</p> : null}
                </div>
                <span aria-hidden="true" className="grid size-11 shrink-0 place-items-center rounded-full bg-brand-50 text-xl text-brand-800">✦</span>
              </div>
              {place.description ? <p className="mt-4 text-sm leading-7 text-stone-600">{place.description}</p> : null}
              {place.directions ? <p className="mt-4 text-sm leading-6 text-stone-700"><strong>Как добраться. </strong>{place.directions}</p> : null}
              <dl className="mt-5 space-y-3 border-t border-stone-100 pt-4 text-sm">
                {place.address ? <div><dt className="font-medium text-stone-900">Адрес</dt><dd className="mt-0.5 text-stone-600"><address className="not-italic">{place.address}</address></dd></div> : null}
                {coordinates ? <div><dt className="font-medium text-stone-900">Координаты</dt><dd className="mt-0.5 text-stone-600">{coordinates}</dd></div> : null}
                {place.phone ? <div><dt className="font-medium text-stone-900">{place.phoneLabel || "Телефон"}</dt><dd className="mt-0.5"><ContentLink href={russianPhone(place.phone).href} className="text-brand-800 underline underline-offset-2">{russianPhone(place.phone).label}</ContentLink></dd></div> : null}
                {place.email ? <div><dt className="font-medium text-stone-900">Почта</dt><dd className="mt-0.5"><ContentLink href={`mailto:${place.email}`} className="text-brand-800 underline underline-offset-2">{place.email}</ContentLink></dd></div> : null}
              </dl>
              <div className="mt-auto flex flex-wrap gap-2 pt-5">
                {mapUrl ? <ContentLink href={mapUrl} className="btn-secondary text-sm">На карте ↗</ContentLink> : null}
                {place.websiteUrl ? <ContentLink href={place.websiteUrl} className="btn-secondary text-sm">Официальный сайт ↗</ContentLink> : null}
                {place.bookingHref ? <ContentLink href={place.bookingHref} className="self-center px-2 text-sm font-medium text-brand-800 underline underline-offset-4">{place.bookingLabel || "Выбрать поездку"}</ContentLink> : null}
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}

function ServiceProfile({ content }: { content: ServiceProfileBlockContent }) {
  const paragraphs = (content.paragraphs ?? []).filter((paragraph) => paragraph.trim());
  const links = (content.links ?? []).filter((link) => link.label.trim() && safeHref(link.href));
  if (paragraphs.length === 0 && links.length === 0) return null;

  return (
    <section className="scroll-mt-24 rounded-3xl bg-brand-50/80 p-6 sm:p-9" id="service">
      <SectionHeading eyebrow={content.eyebrow} title={content.title || "О паломнической службе"} />
      <div className="mt-6 grid gap-8 lg:grid-cols-[1.4fr_0.8fr]">
        <div className="space-y-4 text-sm leading-7 text-stone-700 sm:text-base">
          {paragraphs.map((paragraph, index) => <p key={index}>{paragraph}</p>)}
        </div>
        {links.length ? (
          <div className="space-y-3">
            <p className="text-sm font-semibold uppercase tracking-wide text-stone-500">Полезные сайты</p>
            {links.map((link, index) => (
              <ContentLink key={`${link.href}-${index}`} href={link.href} className="flex items-center justify-between gap-3 rounded-2xl border border-white bg-white px-4 py-4 text-sm font-semibold text-brand-900 transition hover:border-brand-200 hover:shadow-sm">
                <span>{link.label}</span><span aria-hidden="true">↗</span>
              </ContentLink>
            ))}
          </div>
        ) : null}
      </div>
    </section>
  );
}

function ServiceTeam({ content }: { content: ServiceTeamBlockContent }) {
  const members = (content.members ?? []).filter((member) => member.name.trim() || member.role.trim());
  if (members.length === 0) return null;

  return (
    <section className="scroll-mt-24 space-y-6" id="team">
      <SectionHeading eyebrow={content.eyebrow} title={content.title || "Команда службы"} />
      {content.description ? <p className="max-w-3xl text-sm leading-7 text-stone-600">{content.description}</p> : null}
      <div className="relative mx-auto max-w-5xl">
        <div aria-hidden="true" className="absolute left-1/2 top-0 hidden h-6 w-px -translate-x-1/2 bg-brand-200 md:block" />
        <div className="grid gap-4 pt-0 md:grid-cols-2 md:pt-7 lg:grid-cols-3">
          {members.map((member, index) => (
            <article key={`${member.name}-${index}`} className={`flex min-w-0 flex-col items-start gap-4 rounded-2xl border bg-white p-5 shadow-sm md:flex-row ${index === 0 ? "border-brand-300 md:col-span-2 md:mx-auto md:w-2/3 lg:col-span-3 lg:w-1/2" : "border-stone-200"}`}>
              {portraitUrl(member.photo) ? (
                <OptimizedImage src={portraitUrl(member.photo)} alt={member.name}
                  width={112} height={144} sizes="(min-width: 768px) 80px, 112px"
                  className="h-36 w-28 shrink-0 rounded-xl object-cover object-top ring-1 ring-stone-200 md:h-28 md:w-20" />
              ) : null}
              <div className="min-w-0 flex-1 [overflow-wrap:anywhere]">
                <h3 className="font-display text-lg font-semibold leading-snug text-stone-900">{member.name || "Имя не указано"}</h3>
                {member.role ? <p className="mt-1 text-sm font-medium text-stone-600">{member.role}</p> : null}
                {member.bio ? <p className="mt-3 text-sm leading-6 text-stone-600">{member.bio}</p> : null}
                {member.phone || member.email ? (
                  <div className="mt-3 flex flex-col items-start gap-1.5 text-sm">
                    {member.phone ? <ContentLink href={russianPhone(member.phone).href} className="text-brand-800 underline underline-offset-2">{russianPhone(member.phone).label}</ContentLink> : null}
                    {member.email ? <ContentLink href={`mailto:${member.email}`} className="text-brand-800 underline underline-offset-2">{member.email}</ContentLink> : null}
                  </div>
                ) : null}
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

function LegalDonation({ content }: { content: LegalDonationBlockContent }) {
  const legalRows = [
    ["Полное наименование", content.legalName],
    ["Юридический адрес", content.legalAddress],
    ["ИНН", content.inn],
    ["КПП", content.kpp],
    ["ОГРН", content.ogrn],
  ].filter((row): row is [string, string] => Boolean(row[1]?.trim()));
  const donationRows = [
    ["Получатель", content.donationRecipient],
    ["Банк", content.donationBankName],
    ["Расчётный счёт", content.donationAccount],
    ["Корреспондентский счёт", content.donationCorrespondentAccount],
    ["БИК", content.donationBik],
    ["Назначение", content.donationPurpose],
  ].filter((row): row is [string, string] => Boolean(row[1]?.trim()));
  const donationUrl = safeHref(content.donationUrl);

  return (
    <section className="scroll-mt-24 space-y-6" id="documents">
      <SectionHeading eyebrow={content.eyebrow} title={content.title || "Документы и реквизиты"} />
      <div className="grid gap-5 md:grid-cols-[repeat(auto-fit,minmax(18rem,1fr))]">
        <article className="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm sm:p-7">
          <h3 className="font-display text-xl font-semibold text-stone-900">Правовая информация</h3>
          {legalRows.length ? (
            <dl className="mt-5 divide-y divide-stone-100">
              {legalRows.map(([label, value]) => <div key={label} className="grid gap-1 py-3 sm:grid-cols-[12rem_1fr]"><dt className="text-sm text-stone-500">{label}</dt><dd className="break-words text-sm font-medium text-stone-800">{value}</dd></div>)}
            </dl>
          ) : <p className="mt-3 text-sm leading-6 text-stone-600">Юридические документы и условия работы службы.</p>}
          <div className="mt-5 flex flex-wrap gap-2">
            <Link href="/legal" className="btn-secondary text-sm">Все документы</Link>
            <Link href="/privacy" className="btn-secondary text-sm">Политика конфиденциальности</Link>
          </div>
        </article>
          <article className="rounded-3xl bg-amber-50 p-6 sm:p-7">
            <p className="text-xs font-semibold uppercase tracking-[0.15em] text-amber-800">Поддержать службу</p>
            <h3 className="mt-2 font-display text-xl font-semibold text-stone-900">Пожертвование</h3>
            {donationRows.length ? <dl className="mt-4 space-y-2 text-sm text-stone-700">{donationRows.map(([label, value]) => <div key={label}><dt className="font-medium">{label}</dt><dd className="break-words text-stone-600">{value}</dd></div>)}</dl> : null}
            {!donationRows.length && !donationUrl ? (
              <p className="mt-3 text-sm leading-6 text-stone-700">
                Уточните у паломнической службы получателя и назначение перевода.{" "}
                <Link href="/support" className="font-medium text-brand-800 underline underline-offset-2">Связаться со службой</Link>
              </p>
            ) : null}
            {donationUrl ? <ContentLink href={donationUrl} className="btn-primary mt-5">Реквизиты для пожертвования ↗</ContentLink> : null}
          </article>
      </div>
    </section>
  );
}

function LaborTrips({ content }: { content: LaborTripsBlockContent }) {
  const items = (content.items ?? []).filter((item) => item.name.trim());
  if (items.length === 0) return null;

  return (
    <section className="scroll-mt-24 space-y-6" id="labor-trips">
      <SectionHeading eyebrow={content.eyebrow} title={content.title || "Трудовые поездки"} />
      {content.description ? <p className="max-w-3xl text-sm leading-7 text-stone-600">{content.description}</p> : null}
      <div className="grid gap-4 md:grid-cols-2">
        {items.map((item, index) => (
          <article key={`${item.name}-${index}`} className="rounded-3xl border border-stone-200 bg-white p-6 shadow-sm">
            <h3 className="font-display text-xl font-semibold text-stone-900">{item.name}</h3>
            <p className="mt-2 text-sm font-medium text-brand-800">{[item.place, item.dates].filter(Boolean).join(" · ")}</p>
            {item.description ? <p className="mt-4 text-sm leading-7 text-stone-600">{item.description}</p> : null}
            {item.bookingHref ? <ContentLink href={item.bookingHref} className="btn-secondary mt-5">Узнать об участии</ContentLink> : null}
          </article>
        ))}
      </div>
    </section>
  );
}
