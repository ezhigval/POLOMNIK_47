import Link from "next/link";
import { BrandMark } from "@/components/brand-mark";
import { DioceseAffiliation } from "@/components/diocese-affiliation";
import { contactEmail, contactPhone, contactPhoneDisplay } from "@/lib/contact";
import { accountNavLinks, footerNavLinks } from "@/lib/site-nav";
import { siteConfig } from "@/lib/site-config";

export function SiteFooter() {
  return (
    <footer className="mt-auto border-t border-stone-200 bg-stone-900 text-stone-300">
      <div className="mx-auto grid max-w-6xl gap-10 px-4 py-12 sm:grid-cols-2 lg:grid-cols-5">
        <div className="lg:col-span-2">
          <div className="flex items-start gap-3">
            <BrandMark className="size-10 bg-brand-700" />
            <div>
              <p className="text-sm font-semibold leading-snug text-white">{siteConfig.fullName}</p>
            </div>
          </div>
          <p className="mt-4 max-w-md text-sm leading-7 text-stone-400">
            Организуем паломнические поездки по России и миру. Забота и сопровождение на пути к святыням.
          </p>
          <DioceseAffiliation
            className="mt-4 max-w-md"
            textClassName="text-sm leading-6 text-stone-500"
            linkClassName="mt-2 inline-flex items-center text-sm text-stone-300 underline decoration-stone-600 underline-offset-2 transition hover:text-white hover:decoration-white"
          />
        </div>

        <div>
          <p className="mb-4 text-sm font-medium text-white">Навигация</p>
          <ul className="space-y-2.5 text-sm">
            {footerNavLinks.map((link) => (
              <li key={link.href}>
                <Link href={link.href} className="transition hover:text-white">
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </div>

        <div>
          <p className="mb-4 text-sm font-medium text-white">Личный кабинет</p>
          <ul className="space-y-2.5 text-sm">
            <li>
              <Link href="/account/login" className="transition hover:text-white">
                Вход
              </Link>
            </li>
            <li>
              <Link href="/account/register" className="transition hover:text-white">
                Регистрация
              </Link>
            </li>
            {accountNavLinks.map((link) => (
              <li key={link.href}>
                <Link href={link.href} className="transition hover:text-white">
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </div>

        <div>
          <p className="mb-4 text-sm font-medium text-white">Контакты</p>
          <ul className="space-y-2.5 text-sm">
            <li>
              <a href={`tel:${contactPhone}`} className="transition hover:text-white">
                {contactPhoneDisplay}
              </a>
            </li>
            <li>
              <a href={`mailto:${contactEmail}`} className="transition hover:text-white">
                {contactEmail}
              </a>
            </li>
            <li className="pt-2">
              <a
                href={siteConfig.parentOrganization.url}
                target="_blank"
                rel="noopener noreferrer"
                className="underline decoration-stone-600 underline-offset-2 transition hover:text-white hover:decoration-white"
              >
                Сайт Тихвинской епархии ↗
              </a>
            </li>
          </ul>
          <p className="mt-4 text-xs text-stone-500">Пн–Пт, 10:00–19:00 (МСК)</p>
        </div>
      </div>

      <div className="border-t border-stone-800 py-5 text-center text-xs text-stone-500">
        © {new Date().getFullYear()} {siteConfig.name}. Все права защищены.
      </div>
    </footer>
  );
}
