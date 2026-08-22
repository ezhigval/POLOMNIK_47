import Link from "next/link";
import { BrandMark } from "@/components/brand-mark";
import { MobileNav } from "@/components/mobile-nav";
import { UserMenu } from "@/components/user-menu";
import { contactPhone, contactPhoneDisplay } from "@/lib/contact";
import { mainNavLinks } from "@/lib/site-nav";
import { siteConfig } from "@/lib/site-config";
import type { User } from "@/lib/api/auth";

type SiteHeaderProps = {
  user: User | null;
};

export function SiteHeader({ user }: SiteHeaderProps) {
  return (
    <header className="sticky top-0 z-50 border-b border-stone-200/80 bg-white/90 backdrop-blur-md">
      <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-3 sm:py-4">
        <Link href="/" className="group flex items-center gap-3" aria-label={siteConfig.fullName}>
          <BrandMark />
          <span className="flex min-w-0 flex-col">
            <span className="text-base font-semibold leading-tight tracking-tight text-stone-900 group-hover:text-brand-800 sm:text-lg">
              {siteConfig.name}
            </span>
            <span className="text-[11px] leading-tight text-stone-500 sm:text-xs">{siteConfig.tagline}</span>
          </span>
        </Link>

        <nav
          className="hidden items-center gap-1 text-sm lg:flex"
          aria-label="Основная навигация"
        >
          {mainNavLinks.map((link) => (
            <Link key={link.href} href={link.href} className="nav-link">
              {link.label}
            </Link>
          ))}
          <UserMenu user={user} />
          <a href={`tel:${contactPhone}`} className="btn-primary ml-1 px-4 py-2 text-sm">
            {contactPhoneDisplay}
          </a>
        </nav>

        <MobileNav user={user} />
      </div>
    </header>
  );
}
