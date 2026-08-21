"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { contactPhone, contactPhoneDisplay } from "@/lib/contact";
import { accountNavLinks, mainNavLinks } from "@/lib/site-nav";
import type { User } from "@/lib/api/auth";

type MobileNavProps = {
  user: User | null;
};

export function MobileNav({ user }: MobileNavProps) {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    document.body.style.overflow = open ? "hidden" : "";
    return () => {
      document.body.style.overflow = "";
    };
  }, [open]);

  return (
    <div className="lg:hidden">
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        className="inline-flex size-10 items-center justify-center rounded-full border border-stone-200 text-stone-700"
        aria-expanded={open}
        aria-label={open ? "Закрыть меню" : "Открыть меню"}
      >
        {open ? (
          <svg viewBox="0 0 24 24" className="size-5" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M6 6l12 12M6 18L18 6" />
          </svg>
        ) : (
          <svg viewBox="0 0 24 24" className="size-5" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M4 7h16M4 12h16M4 17h16" />
          </svg>
        )}
      </button>

      {open ? (
        <div className="fixed inset-0 top-[57px] z-50 bg-white">
          <nav className="flex flex-col gap-1 p-4" aria-label="Мобильная навигация">
            {mainNavLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                onClick={() => setOpen(false)}
                className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
              >
                {link.label}
              </Link>
            ))}
            {!user ? (
              <Link
                href="/support/chat"
                onClick={() => setOpen(false)}
                className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
              >
                Чат поддержки
              </Link>
            ) : null}
            {user ? (
              accountNavLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  onClick={() => setOpen(false)}
                  className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
                >
                  {link.label}
                </Link>
              ))
            ) : (
              <>
                <Link
                  href="/account/login"
                  onClick={() => setOpen(false)}
                  className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
                >
                  Войти
                </Link>
                <Link
                  href="/account/register"
                  onClick={() => setOpen(false)}
                  className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
                >
                  Регистрация
                </Link>
              </>
            )}
            <a
              href={`tel:${contactPhone}`}
              className="btn-primary mt-4 text-center"
              onClick={() => setOpen(false)}
            >
              {contactPhoneDisplay}
            </a>
          </nav>
        </div>
      ) : null}
    </div>
  );
}
