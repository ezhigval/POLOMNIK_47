export type NavLink = {
  href: string;
  label: string;
};

/** Active state for header / mobile nav (hash anchors excluded). */
export function isMainNavLinkActive(href: string, pathname: string): boolean {
  if (href.includes("#")) {
    return false;
  }
  if (href === "/") {
    return pathname === "/";
  }
  if (href === "/search") {
    return pathname === "/search" || pathname.startsWith("/tours/");
  }
  return pathname === href || pathname.startsWith(`${href}/`);
}

/** Основное меню — один источник для header, footer, mobile */
export const mainNavLinks: NavLink[] = [
  { href: "/", label: "Главная" },
  { href: "/search", label: "Туры" },
  { href: "/news", label: "Новости" },
  { href: "/reviews", label: "Отзывы" },
  { href: "/support", label: "Поддержка" },
];

export const footerNavLinks: NavLink[] = [
  { href: "/", label: "Главная" },
  { href: "/search", label: "Туры" },
  { href: "/news", label: "Новости" },
  { href: "/#about", label: "О службе" },
  { href: "/#why-us", label: "Почему мы" },
  { href: "/#how-it-works", label: "Как записаться" },
  { href: "/reviews", label: "Отзывы" },
  { href: "/#faq", label: "Вопросы и ответы" },
  { href: "/support", label: "Справочник поддержки" },
  { href: "/support/chat", label: "Чат поддержки" },
  { href: "/legal", label: "Юридические документы" },
  { href: "/legal/privacy-policy", label: "Политика конфиденциальности" },
  { href: "/legal/terms", label: "Пользовательское соглашение" },
  { href: "/legal/offer", label: "Публичная оферта" },
];

export const accountNavLinks: NavLink[] = [
  { href: "/account", label: "Профиль" },
  { href: "/account/passengers", label: "Пассажиры" },
  { href: "/account/trips", label: "Мои поездки" },
  { href: "/account/favorites", label: "Избранное" },
  { href: "/account/consents", label: "Согласия" },
  { href: "/account/photos", label: "Фотографии" },
  { href: "/support/chat", label: "Чат поддержки" },
];

/** Burger already has «Поддержка»; chat lives on that page, not as a second item. */
export const burgerAccountNavLinks: NavLink[] = accountNavLinks.filter(
  (link) => link.href !== "/support/chat",
);

export const oauthErrorMessages: Record<string, string> = {
  oauth_not_configured: "Вход через соцсеть пока не настроен. Используйте email или телефон.",
  oauth_cancelled: "Вход через соцсеть отменён. Попробуйте снова или войдите по паролю.",
  oauth_failed: "Не удалось авторизоваться через соцсеть. Попробуйте позже.",
  oauth_profile: "Соцсеть не передала профиль. Попробуйте другой способ входа.",
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
