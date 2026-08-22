"use client";

import Link from "next/link";
import { FormEvent, useEffect, useRef, useState } from "react";
import type { SupportMessage, SupportThread } from "@/lib/auth/user-features";
import { FormError } from "@/components/form-error";

type SupportChatProps = {
  initialThread: SupportThread | null;
  isAuthenticated: boolean;
};

export function SupportChat({ initialThread, isAuthenticated }: SupportChatProps) {
  const [thread, setThread] = useState<SupportThread | null>(initialThread);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const bottomRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [thread?.messages.length]);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!message.trim()) {
      return;
    }

    setLoading(true);
    setError(null);
    const response = await fetch("/api/support", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ body: message.trim() }),
    });

    setLoading(false);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось отправить сообщение");
      return;
    }

    const payload = await response.json();
    setThread(payload.data);
    setMessage("");
  }

  if (!isAuthenticated) {
    return (
      <div className="rounded-2xl border border-stone-200 bg-white p-8 text-center shadow-sm">
        <p className="text-stone-600">Чат поддержки доступен после входа в аккаунт.</p>
        <div className="mt-4 flex flex-wrap justify-center gap-3">
          <Link
            href={`/account/login?returnUrl=${encodeURIComponent("/support/chat")}`}
            className="btn-primary"
          >
            Войти
          </Link>
          <Link href="/support" className="btn-secondary">
            Справочник
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-[min(70vh,640px)] flex-col overflow-hidden rounded-2xl border border-stone-200 bg-white shadow-sm">
      <div className="border-b border-stone-100 px-5 py-4">
        <h2 className="font-semibold text-stone-900">Чат с поддержкой</h2>
        <p className="text-sm text-stone-500">Менеджер ответит в рабочее время.</p>
      </div>

      <div className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
        {(thread?.messages ?? []).length === 0 ? (
          <p className="text-sm text-stone-500">Напишите вопрос — мы на связи.</p>
        ) : (
          thread?.messages.map((item) => <ChatBubble key={item.id} message={item} />)
        )}
        <div ref={bottomRef} />
      </div>

      <form onSubmit={onSubmit} className="border-t border-stone-100 p-4">
        <div className="flex gap-2">
          <input
            value={message}
            onChange={(event) => setMessage(event.target.value)}
            placeholder="Опишите вопрос…"
            className="input-field flex-1"
            disabled={loading}
          />
          <button type="submit" disabled={loading || !message.trim()} className="btn-primary px-5">
            {loading ? "…" : "Отправить"}
          </button>
        </div>
        <FormError className="mt-2">{error}</FormError>
      </form>
    </div>
  );
}

function ChatBubble({ message }: { message: SupportMessage }) {
  const isUser = message.sender_type === "user";
  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`max-w-[85%] rounded-2xl px-4 py-3 text-sm leading-6 ${
          isUser ? "bg-brand-800 text-white" : "bg-stone-100 text-stone-800"
        }`}
      >
        {!isUser ? <p className="mb-1 text-xs font-medium text-stone-500">Поддержка</p> : null}
        {message.body}
      </div>
    </div>
  );
}
