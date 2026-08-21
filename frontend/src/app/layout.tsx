import type { Metadata } from "next";
import { Analytics } from "@/components/analytics";
import { displaySerif, geistMono, geistSans } from "@/lib/fonts";
import { siteConfig } from "@/lib/site-config";
import "./globals.css";

export const metadata: Metadata = {
  metadataBase: new URL(siteConfig.url),
  title: {
    default: `${siteConfig.name} — Паломнические туры по России`,
    template: `%s | ${siteConfig.name}`,
  },
  description: siteConfig.description,
  openGraph: {
    title: `${siteConfig.name} — Паломнические туры`,
    description: siteConfig.description,
    locale: "ru_RU",
    type: "website",
    siteName: siteConfig.name,
  },
  twitter: {
    card: "summary_large_image",
    title: siteConfig.name,
    description: siteConfig.description,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="ru"
      className={`${geistSans.variable} ${geistMono.variable} ${displaySerif.variable} h-full`}
    >
      <body className="flex min-h-full flex-col bg-stone-50 text-stone-900 antialiased">
        {children}
        <Analytics />
      </body>
    </html>
  );
}
