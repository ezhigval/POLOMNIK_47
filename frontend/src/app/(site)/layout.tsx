import { NavigationTracker } from "@/components/navigation-tracker";
import { SiteFooter } from "@/components/site-footer";
import { SiteHeader } from "@/components/site-header";
import { StructuredData } from "@/components/structured-data";
import { getSiteLogo } from "@/lib/site-logo";
import { getSessionUser } from "@/lib/auth/session";

export default async function SiteLayout({ children }: { children: React.ReactNode }) {
  const [user, logo] = await Promise.all([getSessionUser(), getSiteLogo()]);

  return (
    <>
      <NavigationTracker />
      <StructuredData logo={logo} />
      <SiteHeader user={user} />
      <main className="min-w-0 flex-1">{children}</main>
      <SiteFooter />
    </>
  );
}
