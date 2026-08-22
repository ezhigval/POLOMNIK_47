"use client";

import { useEffect, useState } from "react";

type MethodStatus = { available: boolean; message?: string; username?: string };

type AuthMethods = {
  yandex?: MethodStatus;
  vk?: MethodStatus;
  max?: MethodStatus;
  telegram?: MethodStatus;
  mail?: MethodStatus;
};

export function SocialAuthButtons() {
  const [methods, setMethods] = useState<AuthMethods | null>(null);

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

  if (!methods) {
    return null;
  }

  const items: Array<{ key: keyof AuthMethods; label: string; href: string }> = [
    { key: "yandex", label: "Яндекс", href: "/api/auth/social/yandex" },
    { key: "vk", label: "VK", href: "/api/auth/social/vk" },
    { key: "max", label: "Max", href: "/api/auth/social/max" },
  ];

  return (
    <div className="space-y-3">
      <p className="text-center text-xs uppercase tracking-wide text-stone-500">Или войти через</p>
      <div className="grid gap-2 sm:grid-cols-3">
        {items.map((item) => {
          const status = methods[item.key];
          const available = Boolean(status?.available);
          if (!available) {
            return (
              <button
                key={item.key}
                type="button"
                disabled
                className="rounded-xl border border-stone-200 bg-stone-50 px-3 py-2 text-sm text-stone-400"
                title={status?.message || "Пока что недоступно, используйте другой вариант."}
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
