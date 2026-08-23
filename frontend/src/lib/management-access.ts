import type { ManagementSession } from "@/lib/api/management";

export const PERM = {
  tours: "manage_tours",
  bookings: "manage_bookings",
  support: "manage_support",
  content: "manage_content",
  settingsSite: "manage_settings_site",
  recipients: "manage_recipients",
  roles: "manage_roles",
  stats: "view_stats",
  integrations: "manage_integrations",
} as const;

export type ManagementNavItem = {
  href: string;
  label: string;
  exact?: boolean;
  anyOf: string[];
};

export const MANAGEMENT_NAV: ManagementNavItem[] = [
  { href: "/management", label: "Обзор", exact: true, anyOf: [] },
  { href: "/management/content", label: "Главная", anyOf: [PERM.content] },
  { href: "/management/news", label: "Новости", anyOf: [PERM.content] },
  { href: "/management/smm", label: "Контент-план", anyOf: [PERM.content] },
  { href: "/management/tours", label: "Туры", anyOf: [PERM.tours] },
  { href: "/management/bookings", label: "Заявки", anyOf: [PERM.bookings] },
  { href: "/management/support", label: "Поддержка", anyOf: [PERM.support] },
  { href: "/management/reviews", label: "Отзывы", anyOf: [PERM.content] },
  { href: "/management/integrations", label: "Синхронизация", anyOf: [PERM.integrations, PERM.stats] },
  { href: "/management/ai", label: "ИИ и watchdog", anyOf: [PERM.stats] },
  {
    href: "/management/settings",
    label: "Настройки",
    anyOf: [PERM.settingsSite, PERM.recipients, PERM.roles],
  },
];

export function sessionHasPermission(
  session: Pick<ManagementSession, "full_admin" | "permissions"> | null | undefined,
  permission: string,
): boolean {
  if (!session) {
    return false;
  }
  if (session.full_admin) {
    return true;
  }
  return (session.permissions ?? []).includes(permission);
}

export function sessionHasAnyPermission(
  session: Pick<ManagementSession, "full_admin" | "permissions"> | null | undefined,
  permissions: string[],
): boolean {
  if (permissions.length === 0) {
    return Boolean(session);
  }
  return permissions.some((permission) => sessionHasPermission(session, permission));
}

export function visibleManagementNav(
  session: Pick<ManagementSession, "full_admin" | "permissions"> | null | undefined,
): ManagementNavItem[] {
  return MANAGEMENT_NAV.filter((item) => sessionHasAnyPermission(session, item.anyOf));
}
