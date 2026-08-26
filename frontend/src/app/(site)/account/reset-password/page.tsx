import type { Metadata } from "next";
import { PageIntro } from "@/components/page-intro";
import { ResetPasswordForm } from "@/components/auth/reset-password-form";

export const metadata: Metadata = {
  title: "Новый пароль",
  robots: { index: false, follow: false },
};

type PageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function ResetPasswordPage({ searchParams }: PageProps) {
  const params = await searchParams;
  const token = typeof params.token === "string" ? params.token : "";

  return (
    <div className="mx-auto flex min-h-[60vh] max-w-md flex-col justify-center gap-4 px-4 py-12">
      <PageIntro
        backHref="/account/login"
        backLabel="К входу"
        title="Новый пароль"
        description="Задайте пароль по ссылке из письма."
      />
      <ResetPasswordForm token={token} />
    </div>
  );
}
