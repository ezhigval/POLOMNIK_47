import type { Metadata } from "next";
import { Analytics } from "@/components/analytics";
import { CookieBanner } from "@/components/cookie-banner";
import { LiveRefresh } from "@/components/live-refresh";
import { displaySerif, geistMono, geistSans } from "@/lib/fonts";
import { siteConfig } from "@/lib/site-config";
import "./globals.css";

export const viewport = {
  width: "device-width",
  initialScale: 1,
};

export const metadata: Metadata = {
  metadataBase: new URL(siteConfig.url),
  title: {
    default: `${siteConfig.name} — ${siteConfig.tagline}`,
    template: `%s | ${siteConfig.name}`,
  },
  description: siteConfig.description,
  openGraph: {
    title: `${siteConfig.name} — ${siteConfig.tagline}`,
    description: siteConfig.description,
    locale: "ru_RU",
    type: "website",
    siteName: siteConfig.name,
    url: siteConfig.url,
  },
  twitter: {
    card: "summary_large_image",
    title: `${siteConfig.name} — ${siteConfig.tagline}`,
    description: siteConfig.description,
  },
  robots: {
    index: true,
    follow: true,
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
        {process.env.NEXT_PUBLIC_LIVE_REFRESH === "1" ? <LiveRefresh /> : null}
        <CookieBanner />
        <Analytics />
      </body>
    </html>
  );
}
