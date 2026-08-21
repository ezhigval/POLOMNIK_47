import Link from "next/link";
import { contactEmail, contactPhone, contactPhoneDisplay } from "@/lib/contact";
import { accountNavLinks, footerNavLinks } from "@/lib/site-nav";
import { siteConfig } from "@/lib/site-config";

export function SiteFooter() {
  return (
    <footer className="mt-auto border-t border-stone-200 bg-stone-900 text-stone-300">
      <div className="mx-auto grid max-w-6xl gap-10 px-4 py-12 sm:grid-cols-2 lg:grid-cols-5">
        <div className="lg:col-span-2">
          <div className="flex items-center gap-3">
            <span className="flex size-10 items-center justify-center rounded-full bg-brand-700 font-display text-sm font-bold text-white">
              47
            </span>
            <div>
              <p className="font-semibold text-white">{siteConfig.name}</p>
              <p className="text-sm text-stone-400">{siteConfig.tagline}</p>
            </div>
          </div>
          <p className="mt-4 max-w-md text-sm leading-7 text-stone-400">
            Организуем паломнические поездки по России: от Оптиной пустыни до Валаама. Забота,
            сопровождение и тишина на пути к святыням.
          </p>
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
