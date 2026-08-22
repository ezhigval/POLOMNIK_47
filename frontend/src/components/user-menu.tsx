"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import type { User } from "@/lib/api/auth";
import { accountNavLinks, userContactLine } from "@/lib/site-nav";

type UserMenuProps = {
  user: User | null;
};

export function UserMenu({ user }: UserMenuProps) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  async function logout() {
    setLoading(true);
    await fetch("/api/auth/logout", { method: "POST" });
    setLoading(false);
    setOpen(false);
    router.refresh();
    router.push("/");
  }

  if (!user) {
    return (
      <div className="hidden items-center gap-2 lg:flex">
        <Link href="/account/login" className="btn-secondary px-4 py-2 text-sm">
          Войти
        </Link>
        <Link href="/account/register" className="btn-primary px-4 py-2 text-sm">
          Регистрация
        </Link>
      </div>
    );
  }

  const initials = user.name
    .split(" ")
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("");

  return (
    <div className="relative hidden lg:block">
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        className="inline-flex items-center gap-2 rounded-full border border-stone-200 bg-white px-3 py-2 text-sm font-medium text-stone-800 transition hover:bg-stone-50"
      >
        <span className="flex size-8 items-center justify-center rounded-full bg-brand-800 text-xs font-semibold text-white">
          {initials || "Я"}
        </span>
        <span className="max-w-[120px] truncate">{user.name.split(" ")[0]}</span>
      </button>

      {open ? (
        <>
          <button
            type="button"
            className="fixed inset-0 z-40 cursor-default"
            aria-label="Закрыть меню"
            onClick={() => setOpen(false)}
          />
          <div className="absolute right-0 z-50 mt-2 w-56 overflow-hidden rounded-2xl border border-stone-200 bg-white py-2 shadow-lg">
            <div className="border-b border-stone-100 px-4 py-3">
              <p className="font-medium text-stone-900">{user.name}</p>
              <p className="truncate text-xs text-stone-500">{userContactLine(user)}</p>
            </div>
            {accountNavLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="block px-4 py-2.5 text-sm text-stone-700 hover:bg-stone-50"
                onClick={() => setOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            <Link
              href="/support"
              className="block px-4 py-2.5 text-sm text-stone-700 hover:bg-stone-50"
              onClick={() => setOpen(false)}
            >
              Справочник
            </Link>
            <button
              type="button"
              onClick={logout}
              disabled={loading}
              className="block w-full px-4 py-2.5 text-left text-sm text-red-700 hover:bg-red-50 disabled:opacity-60"
            >
              {loading ? "Выход…" : "Выйти"}
            </button>
          </div>
        </>
      ) : null}
    </div>
  );
}
