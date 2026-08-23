"use client";

import { useEffect, useState } from "react";
import { PersonalDataConsentCheckbox, TermsConsentCheckbox } from "@/components/consent-checkbox";

type MethodStatus = { available: boolean; message?: string; username?: string };

type AuthMethods = {
  yandex?: MethodStatus;
  vk?: MethodStatus;
  max?: MethodStatus;
  telegram?: MethodStatus;
  mail?: MethodStatus;
};

const OAUTH_CONSENT_COOKIE = "oauth_legal_consent";

function setOAuthConsentCookie(personalData: boolean, terms: boolean) {
  if (!personalData || !terms) {
    document.cookie = `${OAUTH_CONSENT_COOKIE}=; Path=/; Max-Age=0; SameSite=Lax`;
    return;
  }
  document.cookie = `${OAUTH_CONSENT_COOKIE}=pd,terms; Path=/; Max-Age=600; SameSite=Lax`;
}

type SocialAuthButtonsProps = {
  /** Если true — для первого входа нужны согласия (регистрация / новый аккаунт). */
  requireConsentForNewAccount?: boolean;
};

export function SocialAuthButtons({ requireConsentForNewAccount = true }: SocialAuthButtonsProps) {
  const [methods, setMethods] = useState<AuthMethods | null>(null);
  const [consentPersonalData, setConsentPersonalData] = useState(false);
  const [consentTerms, setConsentTerms] = useState(false);

  useEffect(() => {
    let cancelled = false;
    fetch("/api/v1/auth/methods", { cache: "no-store" })
      .then((r) => r.json())
      .then((body) => {
        if (!cancelled) setMethods(body?.data ?? null);
      })
      .catch(() => {
        if (!cancelled) setMethods({});
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (!requireConsentForNewAccount) return;
    setOAuthConsentCookie(consentPersonalData, consentTerms);
  }, [consentPersonalData, consentTerms, requireConsentForNewAccount]);

  if (!methods) {
    return null;
  }

  const items: Array<{ key: keyof AuthMethods; label: string; href: string }> = [
    { key: "yandex", label: "Яндекс", href: "/api/auth/social/yandex" },
    { key: "vk", label: "VK", href: "/api/auth/social/vk" },
    { key: "max", label: "Max", href: "/api/auth/social/max" },
  ];

  const canStartOAuth = !requireConsentForNewAccount || (consentPersonalData && consentTerms);

  return (
    <div className="space-y-3">
      <p className="text-center text-xs uppercase tracking-wide text-stone-500">Или войти через</p>
      {requireConsentForNewAccount ? (
        <div className="space-y-3 rounded-xl border border-stone-100 bg-stone-50 p-3">
          <p className="text-xs text-stone-600">
            При создании нового аккаунта через соцсеть нужны согласия. Для уже существующего аккаунта
            они не требуются повторно.
          </p>
          <PersonalDataConsentCheckbox checked={consentPersonalData} onChange={setConsentPersonalData} />
          <TermsConsentCheckbox checked={consentTerms} onChange={setConsentTerms} />
        </div>
      ) : null}
      <div className="grid gap-2 sm:grid-cols-3">
        {items.map((item) => {
          const status = methods[item.key];
          const available = Boolean(status?.available);
          if (!available || !canStartOAuth) {
            return (
              <button
                key={item.key}
                type="button"
                disabled
                className="rounded-xl border border-stone-200 bg-stone-50 px-3 py-2 text-sm text-stone-400"
                title={
                  !available
                    ? status?.message || "Пока что недоступно, используйте другой вариант."
                    : "Отметьте согласия для первого входа через соцсеть"
                }
              >
                {item.label}
              </button>
            );
          }
          return (
            <a key={item.key} href={item.href} className="btn-secondary text-center text-sm">
              {item.label}
            </a>
          );
        })}
      </div>
      {methods.telegram?.available && methods.telegram.username ? (
        <p className="text-center text-xs text-stone-500">
          Telegram Login: виджет для @{methods.telegram.username} (домен в BotFather).
        </p>
      ) : (
        <p className="text-center text-xs text-stone-400">
          Telegram Login: {methods.telegram?.message || "пока недоступно"}
        </p>
      )}
    </div>
  );
}
