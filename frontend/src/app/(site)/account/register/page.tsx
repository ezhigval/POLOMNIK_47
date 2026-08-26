import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { PageIntro } from "@/components/page-intro";
import { RegisterForm } from "@/components/auth/register-form";
import { getSessionUser } from "@/lib/auth/session";
import { safeReturnUrl } from "@/lib/site-nav";

export const metadata: Metadata = {
  title: "Регистрация",
  robots: { index: false, follow: false },
};

type RegisterPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function RegisterPage({ searchParams }: RegisterPageProps) {
  const user = await getSessionUser();
  const params = await searchParams;
  const returnUrl = safeReturnUrl(typeof params.returnUrl === "string" ? params.returnUrl : undefined);

  if (user) {
    redirect(returnUrl);
  }

  return (
    <div className="mx-auto flex min-h-[60vh] max-w-md flex-col justify-center gap-4 px-4 py-12">
      <PageIntro
        backHref="/account/login"
        backLabel="К входу"
        title="Регистрация"
        description="Создайте аккаунт, чтобы сохранять туры и отслеживать заявки."
      />
      <RegisterForm returnUrl={returnUrl} />
    </div>
  );
}
