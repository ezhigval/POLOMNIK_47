export type NavLink = {
  href: string;
  label: string;
};

/** Основное меню — один источник для header, footer, mobile */
export const mainNavLinks: NavLink[] = [
  { href: "/search", label: "Поиск туров" },
  { href: "/#about", label: "О службе" },
  { href: "/#how-it-works", label: "Как записаться" },
  { href: "/reviews", label: "Отзывы" },
  { href: "/support", label: "Поддержка" },
];

export const footerNavLinks: NavLink[] = [
  { href: "/search", label: "Поиск туров" },
  { href: "/#about", label: "О службе" },
  { href: "/#why-us", label: "Почему мы" },
  { href: "/#how-it-works", label: "Как записаться" },
  { href: "/reviews", label: "Отзывы" },
  { href: "/#faq", label: "Вопросы и ответы" },
  { href: "/support", label: "Справочник поддержки" },
  { href: "/support/chat", label: "Чат поддержки" },
  { href: "/privacy", label: "Политика конфиденциальности" },
];

export const accountNavLinks: NavLink[] = [
  { href: "/account/trips", label: "Мои поездки" },
  { href: "/account/favorites", label: "Избранное" },
  { href: "/support/chat", label: "Чат поддержки" },
];

export const oauthErrorMessages: Record<string, string> = {
  oauth_not_configured: "Вход через Google пока не настроен. Используйте email или телефон.",
  oauth_cancelled: "Вход через Google отменён. Попробуйте снова или войдите по паролю.",
  oauth_failed: "Не удалось авторизоваться через Google. Попробуйте позже.",
  oauth_profile: "Google не передал профиль. Попробуйте другой способ входа.",
  oauth_backend: "Сервер не смог создать аккаунт. Попробуйте регистрацию вручную.",
};

export function safeReturnUrl(value: string | null | undefined): string {
  if (!value || !value.startsWith("/") || value.startsWith("//")) {
    return "/account/trips";
  }
  return value;
}

export function userContactLine(user: { phone?: string; email?: string }): string {
  if (user.phone?.trim()) {
    return user.phone;
  }
  if (user.email?.trim()) {
    return user.email;
  }
  return "Профиль без контактов";
}
