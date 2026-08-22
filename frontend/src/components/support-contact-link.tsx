"use client";

import type { ComponentProps, ReactNode } from "react";
import Link from "next/link";
import { trackSupportContact } from "@/lib/analytics";

type Channel = "phone" | "email" | "chat";

type SupportContactLinkProps = {
  channel: Channel;
  children: ReactNode;
  className?: string;
} & (
  | { href: string; as?: "a" }
  | { href: ComponentProps<typeof Link>["href"]; as: "link" }
);

export function SupportContactLink({
  channel,
  href,
  children,
  className,
  as = "a",
}: SupportContactLinkProps) {
  function onClick() {
    trackSupportContact(channel);
  }

  if (as === "link") {
    return (
      <Link href={href} className={className} onClick={onClick}>
        {children}
      </Link>
    );
  }

  return (
    <a href={String(href)} className={className} onClick={onClick}>
      {children}
    </a>
  );
}
