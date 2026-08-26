"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useId, useState, useSyncExternalStore } from "react";
import { createPortal } from "react-dom";
import { desktopMediaQuery } from "@/lib/breakpoints";
import { contactPhone, contactPhoneDisplay } from "@/lib/contact";
import { accountNavLinks, isMainNavLinkActive, mainNavLinks } from "@/lib/site-nav";
import type { User } from "@/lib/api/auth";

type MobileNavProps = {
  user: User | null;
};

function subscribeNoop() {
  return () => {};
}

export function MobileNav({ user }: MobileNavProps) {
  const pathname = usePathname();
  const [openPath, setOpenPath] = useState<string | null>(null);
  const open = openPath === pathname;
  const isClient = useSyncExternalStore(subscribeNoop, () => true, () => false);
  const titleId = useId();

  useEffect(() => {
    const media = window.matchMedia(desktopMediaQuery);
    function onChange() {
      if (media.matches) {
        setOpenPath(null);
      }
    }

    media.addEventListener("change", onChange);
    return () => media.removeEventListener("change", onChange);
  }, []);

  useEffect(() => {
    if (!open) {
      return;
    }

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setOpenPath(null);
      }
    }

    window.addEventListener("keydown", onKeyDown);
    return () => {
      document.body.style.overflow = previousOverflow;
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  function close() {
    setOpenPath(null);
  }

  const menu =
    open && isClient
      ? createPortal(
          <div className="fixed inset-0 z-[100] lg:hidden">
            <button
              type="button"
              className="absolute inset-0 bg-black/50"
              aria-label="Закрыть меню"
              onClick={close}
            />
            <nav
              role="dialog"
              aria-modal="true"
              aria-labelledby={titleId}
              className="absolute inset-y-0 right-0 flex w-[min(20rem,calc(100%-3.5rem))] flex-col bg-white pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)] shadow-2xl"
            >
              <div className="flex items-center justify-between border-b border-stone-100 px-4 py-3">
                <p id={titleId} className="text-sm font-semibold text-stone-900">
                  Меню
                </p>
                <button
                  type="button"
                  onClick={close}
                  className="inline-flex size-10 items-center justify-center rounded-full border border-stone-200 text-stone-700"
                  aria-label="Закрыть меню"
                >
                  <svg viewBox="0 0 24 24" className="size-5" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M6 6l12 12M6 18L18 6" />
                  </svg>
                </button>
              </div>

              <div className="flex flex-1 flex-col gap-1 overflow-y-auto p-3">
                {mainNavLinks.map((link) => {
                  const active = isMainNavLinkActive(link.href, pathname);
                  return (
                    <Link
                      key={link.href}
                      href={link.href}
                      onClick={close}
                      aria-current={active ? "page" : undefined}
                      className={`rounded-xl px-4 py-3 text-base font-medium transition hover:bg-stone-50 ${
                        active ? "bg-brand-50 text-brand-900" : "text-stone-800"
                      }`}
                    >
                      {link.label}
                    </Link>
                  );
                })}
                {!user ? (
                  <Link
                    href="/support/chat"
                    onClick={close}
                    className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
                  >
                    Чат поддержки
                  </Link>
                ) : null}
                {user ? (
                  accountNavLinks.map((link) => {
                    const active =
                      link.href === "/account"
                        ? pathname === "/account"
                        : pathname === link.href || pathname.startsWith(`${link.href}/`);
                    return (
                      <Link
                        key={link.href}
                        href={link.href}
                        onClick={close}
                        aria-current={active ? "page" : undefined}
                        className={`rounded-xl px-4 py-3 text-base font-medium transition hover:bg-stone-50 ${
                          active ? "bg-brand-50 text-brand-900" : "text-stone-800"
                        }`}
                      >
                        {link.label}
                      </Link>
                    );
                  })
                ) : (
                  <>
                    <Link
                      href="/account/login"
                      onClick={close}
                      className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
                    >
                      Войти
                    </Link>
                    <Link
                      href="/account/register"
                      onClick={close}
                      className="rounded-xl px-4 py-3 text-base font-medium text-stone-800 hover:bg-stone-50"
                    >
                      Регистрация
                    </Link>
                  </>
                )}
              </div>

              <div className="border-t border-stone-100 p-4">
                <a href={`tel:${contactPhone}`} className="btn-primary w-full text-center" onClick={close}>
                  {contactPhoneDisplay}
                </a>
              </div>
            </nav>
          </div>,
          document.body,
        )
      : null;

  return (
    <div className="lg:hidden">
      <button
        type="button"
        onClick={() => setOpenPath(open ? null : pathname)}
        className="inline-flex size-10 items-center justify-center rounded-full border border-stone-200 bg-white text-stone-700"
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
      {menu}
    </div>
  );
}
