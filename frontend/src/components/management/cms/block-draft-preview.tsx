"use client";

import Image from "next/image";
import type { CmsBlockType } from "@/lib/api/cms";
import { CtaSection } from "@/components/cta-section";
import { HeroSection } from "@/components/hero-section";
import { SectionHeading } from "@/components/section-heading";

export const livePreviewTypes = new Set<CmsBlockType>([
  "hero", "information_hero", "cta", "rich_text",
]);

function text(content: Record<string, unknown>, key: string): string {
  return typeof content[key] === "string" ? content[key] : "";
}

/**
 * Instant client-only preview of selected editable blocks.
 * It never writes draft data to the server or publishes the CMS page.
 */
export function BlockDraftPreview({
  type,
  content,
  visible,
}: {
  type: CmsBlockType;
  content: Record<string, unknown>;
  visible: boolean;
}) {
  return (
    <aside aria-label="Предпросмотр блока до сохранения" className="min-w-0 space-y-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h3 className="text-sm font-semibold text-stone-900">Предпросмотр при редактировании</h3>
        <span className="rounded-full bg-amber-100 px-2 py-1 text-xs font-medium text-amber-900">
          Не сохранено
        </span>
      </div>
      <p className="text-xs leading-5 text-stone-500">
        Изменения появляются здесь сразу и не видны посетителям сайта.
      </p>
      {!visible ? (
        <p role="status" className="rounded-xl border border-amber-200 bg-amber-50 p-3 text-sm text-amber-900">
          Блок скрыт на сайте. Предпросмотр доступен только в редакторе.
        </p>
      ) : null}
      <div className="overflow-hidden rounded-xl border border-stone-200 bg-white p-2 shadow-sm sm:p-3">
        {type === "hero" ? (
          <HeroSection content={{
            eyebrow: text(content, "eyebrow"),
            title: text(content, "title"),
            subtitle: text(content, "subtitle"),
          }} />
        ) : null}
        {type === "information_hero" ? (
          <section className="relative min-h-72 overflow-hidden rounded-2xl">
            <Image src="/images/hero/tikhvin-monastery.webp" alt="" fill
              sizes="(max-width: 768px) 100vw, 550px" className="object-cover" />
            <div className="absolute inset-0 bg-gradient-to-r from-black/80 via-black/50 to-black/10" />
            <div className="relative flex min-h-72 items-end p-5 text-white sm:p-7">
              <div className="space-y-3">
                {text(content, "eyebrow") ? (
                  <p className="text-xs font-semibold uppercase tracking-widest text-amber-100">
                    {text(content, "eyebrow")}
                  </p>
                ) : null}
                <h2 className="font-display text-3xl font-semibold leading-tight sm:text-4xl">
                  {text(content, "title") || "Заголовок баннера"}
                </h2>
                {text(content, "subtitle") ? (
                  <p className="whitespace-pre-wrap text-sm leading-6 text-white/90">
                    {text(content, "subtitle")}
                  </p>
                ) : null}
                <div className="flex flex-wrap gap-2 pt-2 text-xs">
                  {text(content, "primaryCta") ? (
                    <span className="rounded-full bg-amber-300 px-4 py-2 font-medium text-stone-950">{text(content, "primaryCta")}</span>
                  ) : null}
                  {text(content, "secondaryCta") ? (
                    <span className="rounded-full border border-white/60 px-4 py-2 font-medium">{text(content, "secondaryCta")}</span>
                  ) : null}
                </div>
              </div>
            </div>
          </section>
        ) : null}
        {type === "cta" ? (
          <CtaSection content={{
            title: text(content, "title"),
            subtitle: text(content, "subtitle"),
            button: text(content, "button"),
            href: text(content, "href") || "/search",
          }} />
        ) : null}
        {type === "rich_text" ? (
          <section className="space-y-4 p-3 sm:p-5">
            <SectionHeading eyebrow={text(content, "eyebrow")}
              title={text(content, "title") || "Текст"} />
            <p className="whitespace-pre-wrap break-words text-sm leading-7 text-stone-700">
              {text(content, "body") || "Текст пока не заполнен"}
            </p>
          </section>
        ) : null}
      </div>
      <p className="text-xs leading-5 text-stone-500">
        Компактный предпросмотр одного блока. Полную страницу в выбранной ширине экрана
        смотрите через «Предпросмотр страницы» после сохранения.
      </p>
    </aside>
  );
}
