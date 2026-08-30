import type { ReactNode, SyntheticEvent } from "react";

type AuthExpandableProps = {
  title: string;
  hint?: string;
  children: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
};

export function AuthExpandable({ title, hint, children, open, onOpenChange }: AuthExpandableProps) {
  const controlled = open !== undefined;

  function onToggle(event: SyntheticEvent<HTMLDetailsElement>) {
    onOpenChange?.(event.currentTarget.open);
  }

  return (
    <details
      className="group rounded-xl border border-stone-200 bg-stone-50 open:bg-white"
      {...(controlled ? { open } : {})}
      onToggle={onToggle}
    >
      <summary className="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-sm font-medium text-stone-800 marker:content-none [&::-webkit-details-marker]:hidden">
        <span>
          {title}
          {hint ? <span className="mt-0.5 block text-xs font-normal text-stone-500">{hint}</span> : null}
        </span>
        <span aria-hidden className="text-stone-400 transition group-open:rotate-180">
          ▾
        </span>
      </summary>
      <div className="space-y-3 border-t border-stone-100 px-4 py-3">{children}</div>
    </details>
  );
}
