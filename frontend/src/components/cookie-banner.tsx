"use client";

import Link from "next/link";
import { useEffect, useState, useSyncExternalStore } from "react";
import {
  COOKIE_CONSENT_EVENT,
  COOKIE_CONSENT_KEY,
  type CookieConsentChoice,
  getCookieConsent,
  setCookieConsent,
  subscribeCookieConsent,
} from "@/lib/cookie-consent";

function subscribeNoop() {
  return () => {};
}

export function CookieBanner() {
  const isClient = useSyncExternalStore(subscribeNoop, () => true, () => false);
  const choice = useSyncExternalStore(subscribeCookieConsent, getCookieConsent, () => null);
  const [forceOpen, setForceOpen] = useState(false);

  function choose(next: CookieConsentChoice) {
    setCookieConsent(next);
    setForceOpen(false);
    void fetch("/api/consents", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        consent_type:
          next === "all" ? "cookie_all" : next === "essential" ? "cookie_essential" : "cookie_reject",
      }),
    }).catch(() => {
      /* best-effort server log */
    });
  }

  useEffect(() => {
    function onOpen() {
      setForceOpen(true);
    }
    window.addEventListener(COOKIE_CONSENT_EVENT, onOpen);
    return () => window.removeEventListener(COOKIE_CONSENT_EVENT, onOpen);
  }, []);

  const visible = forceOpen || (isClient && choice === null);
  if (!visible) {
    return null;
  }

  return (
    <div
      role="dialog"
      aria-label="Настройки cookie"
      className="fixed inset-x-0 bottom-0 z-50 border-t border-stone-200 bg-white/95 p-4 shadow-lg backdrop-blur sm:p-5"
    >
      <div className="mx-auto flex max-w-4xl flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="max-w-2xl text-sm leading-6 text-stone-700">
          <p className="font-medium text-stone-900">Файлы cookie</p>
          <p className="mt-1">
            Мы используем необходимые cookie для работы сайта. Аналитические cookie (Яндекс.Метрика,
            Google Analytics) включаются только после вашего согласия. Подробнее — в{" "}
            <Link href="/legal/cookie-policy" className="font-medium text-brand-800 underline underline-offset-2">
              Политике cookie
            </Link>
            .
          </p>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row sm:shrink-0">
          <button type="button" className="btn-secondary text-sm" onClick={() => choose("reject")}>
            Отказаться
          </button>
          <button type="button" className="btn-secondary text-sm" onClick={() => choose("essential")}>
            Только необходимые
          </button>
          <button type="button" className="btn-primary text-sm" onClick={() => choose("all")}>
            Принять все
          </button>
        </div>
      </div>
      <p className="sr-only">Хранение выбора: {COOKIE_CONSENT_KEY}</p>
    </div>
  );
}

export function openCookieSettings() {
  if (typeof window === "undefined") {
    return;
  }
  window.dispatchEvent(new Event(COOKIE_CONSENT_EVENT));
}
