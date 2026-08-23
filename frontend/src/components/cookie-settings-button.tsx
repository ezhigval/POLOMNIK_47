"use client";

import { openCookieSettings } from "@/components/cookie-banner";

export function CookieSettingsButton() {
  return (
    <button
      type="button"
      onClick={() => openCookieSettings()}
      className="text-left text-sm text-stone-400 transition hover:text-white"
    >
      Настройки cookie
    </button>
  );
}
