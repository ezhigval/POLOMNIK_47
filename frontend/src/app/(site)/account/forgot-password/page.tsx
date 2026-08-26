import type { Metadata } from "next";
import { PageIntro } from "@/components/page-intro";
import { ForgotPasswordForm } from "@/components/auth/forgot-password-form";

export const metadata: Metadata = {
  title: "Восстановление пароля",
  robots: { index: false, follow: false },
};

export default function ForgotPasswordPage() {
  return (
    <div className="mx-auto flex min-h-[60vh] max-w-md flex-col justify-center gap-4 px-4 py-12">
      <PageIntro
        backHref="/account/login"
        backLabel="К входу"
        title="Забыли пароль?"
        description="Пришлём ссылку на email, если почта настроена на сервере."
      />
      <ForgotPasswordForm />
    </div>
  );
}
