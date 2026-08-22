"use client";

import { FormEvent, useEffect, useState } from "react";
import { FormError } from "@/components/form-error";

export type PhoneCallChallenge = {
  check_id: string;
  call_phone: string;
  call_phone_pretty: string;
  expires_in: number;
};

type PhoneCallVerifyProps = {
  phone: string;
  disabled?: boolean;
  onConfirmed: (checkId: string) => void;
  onAvailable?: () => void;
  onUnavailable?: (message: string) => void;
};

const UNAVAILABLE = "Пока что недоступно, используйте другой вариант.";

export function PhoneCallVerify({
  phone,
  disabled,
  onConfirmed,
  onAvailable,
  onUnavailable,
}: PhoneCallVerifyProps) {
  const [available, setAvailable] = useState<boolean | null>(null);
  const [message, setMessage] = useState(UNAVAILABLE);
  const [challenge, setChallenge] = useState<PhoneCallChallenge | null>(null);
  const [status, setStatus] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const response = await fetch("/api/v1/auth/methods", { cache: "no-store" });
        const body = await response.json().catch(() => null);
        if (cancelled) return;
        const phoneCall = body?.data?.phone_call;
        const isAvailable = Boolean(phoneCall?.available);
        setAvailable(isAvailable);
        const text = typeof phoneCall?.message === "string" && phoneCall.message ? phoneCall.message : UNAVAILABLE;
        setMessage(text);
        if (isAvailable) {
          onAvailable?.();
        } else {
          onUnavailable?.(text);
        }
      } catch {
        if (!cancelled) {
          setAvailable(false);
          setMessage(UNAVAILABLE);
          onUnavailable?.(UNAVAILABLE);
        }
      }
    })();
    return () => {
      cancelled = true;
    };
    // Intentional: run once on mount for availability.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!challenge || status === "confirmed") return;
    const timer = window.setInterval(async () => {
      const response = await fetch(
        `/api/v1/auth/phone/status?check_id=${encodeURIComponent(challenge.check_id)}`,
        { cache: "no-store" },
      );
      const body = await response.json().catch(() => null);
      const next = body?.data?.status as string | undefined;
      if (!next) return;
      setStatus(next);
      if (next === "confirmed") {
        onConfirmed(challenge.check_id);
      }
      if (next === "expired") {
        setError("Время на звонок истекло. Запросите номер снова.");
        setChallenge(null);
      }
    }, 2500);
    return () => window.clearInterval(timer);
  }, [challenge, status, onConfirmed]);

  async function start(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setLoading(true);
    setStatus(null);
    try {
      const response = await fetch("/api/v1/auth/phone/start", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ phone }),
      });
      const body = await response.json().catch(() => null);
      if (!response.ok) {
        setError(body?.error?.message ?? UNAVAILABLE);
        setLoading(false);
        return;
      }
      setChallenge(body.data as PhoneCallChallenge);
    } catch {
      setError(UNAVAILABLE);
    } finally {
      setLoading(false);
    }
  }

  if (available === null) {
    return <p className="text-sm text-stone-500">Проверяем доступность подтверждения телефона…</p>;
  }

  if (!available) {
    return (
      <p className="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-900">
        Подтверждение телефона звонком: {message}
      </p>
    );
  }

  return (
    <div className="space-y-3 rounded-xl border border-stone-200 bg-stone-50 p-4">
      <p className="text-sm text-stone-700">
        Подтверждение — <strong>звонком с вашего номера</strong> на указанный ниже номер (не SMS-код). Звонок
        сбросят, он бесплатный.
      </p>
      {!challenge ? (
        <button
          type="button"
          disabled={disabled || loading || !phone.trim()}
          onClick={start}
          className="btn-primary w-full"
        >
          {loading ? "Запрашиваем номер…" : "Получить номер для звонка"}
        </button>
      ) : (
        <div className="space-y-2 text-sm text-stone-800">
          <p>
            Позвоните с <span className="font-medium">{phone}</span> на:
          </p>
          <p className="font-display text-xl font-semibold tracking-wide">
            <a className="text-brand-800 underline" href={`tel:${challenge.call_phone}`}>
              {challenge.call_phone_pretty || challenge.call_phone}
            </a>
          </p>
          <p className="text-stone-600">
            {status === "confirmed"
              ? "Номер подтверждён."
              : "Ожидаем звонок (до 5 минут). Страница обновит статус сама."}
          </p>
        </div>
      )}
      <FormError>{error}</FormError>
    </div>
  );
}
