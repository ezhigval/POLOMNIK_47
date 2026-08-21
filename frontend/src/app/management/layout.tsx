import { notFound } from "next/navigation";

export const dynamic = "force-dynamic";

export default function ManagementRootLayout({ children }: { children: React.ReactNode }) {
  if (!process.env.ADMIN_TOKEN) {
    notFound();
  }

  return children;
}
