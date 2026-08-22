import type { Metadata } from "next";
import { notFound } from "next/navigation";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default function ManagementRootLayout({ children }: { children: React.ReactNode }) {
  if (!process.env.ADMIN_TOKEN) {
    notFound();
  }

  return children;
}
