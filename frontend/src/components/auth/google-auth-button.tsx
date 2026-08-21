import Link from "next/link";
import { safeReturnUrl } from "@/lib/site-nav";

type GoogleAuthButtonProps = {
  returnUrl?: string;
};

export function GoogleAuthButton({ returnUrl }: GoogleAuthButtonProps) {
  const destination = safeReturnUrl(returnUrl);
  const href =
    destination === "/account/trips"
      ? "/api/auth/google"
      : `/api/auth/google?returnUrl=${encodeURIComponent(destination)}`;

  return (
    <Link
      href={href}
      className="inline-flex w-full items-center justify-center gap-2 rounded-full border border-stone-300 bg-white px-5 py-2.5 text-sm font-medium text-stone-800 transition hover:bg-stone-50"
    >
      <span className="flex size-5 items-center justify-center rounded-full bg-white text-sm font-bold text-[#4285F4]">
        G
      </span>
      Продолжить через Google
    </Link>
  );
}

export function AuthDivider() {
  return (
    <div className="flex items-center gap-3 py-1">
      <span className="h-px flex-1 bg-stone-200" />
      <span className="text-xs uppercase tracking-wide text-stone-400">или</span>
      <span className="h-px flex-1 bg-stone-200" />
    </div>
  );
}
