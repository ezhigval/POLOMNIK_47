import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { AuthAlert } from "@/components/auth/auth-alert";
import { LoginForm } from "@/components/auth/login-form";
import { PageIntro } from "@/components/page-intro";
import { getSessionUser } from "@/lib/auth/session";
import { oauthErrorMessages, safeReturnUrl } from "@/lib/site-nav";

export const metadata: Metadata = {
  title: "Вход",
  robots: { index: false, follow: false },
};

type LoginPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function LoginPage({ searchParams }: LoginPageProps) {
  const user = await getSessionUser();
  const params = await searchParams;
  const returnUrl = safeReturnUrl(typeof params.returnUrl === "string" ? params.returnUrl : undefined);

  if (user) {
    redirect(returnUrl);
  }

  const errorKey = typeof params.error === "string" ? params.error : undefined;
  const oauthError = errorKey ? oauthErrorMessages[errorKey] : undefined;

  return (
    <div className="mx-auto flex min-h-[60vh] max-w-md flex-col justify-center gap-4 px-4 py-12">
      <PageIntro backHref="/" title="Вход в аккаунт" description="Избранное, заявки и чат поддержки — в одном месте." />
      {oauthError ? <AuthAlert message={oauthError} /> : null}
      <LoginForm returnUrl={returnUrl} />
    </div>
  );
}
