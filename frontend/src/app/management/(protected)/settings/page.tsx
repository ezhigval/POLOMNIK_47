import { ManagementSettingsForms } from "@/components/management/management-settings-forms";
import {
  getManagementNotificationSettings,
  getManagementSession,
  getManagementSiteSettings,
  listManagementRoles,
} from "@/lib/api/management";

const permissionOptions = [
  { id: "manage_tours", label: "Туры" },
  { id: "manage_bookings", label: "Заявки" },
  { id: "manage_content", label: "Контент" },
  { id: "manage_settings_site", label: "Сайт" },
  { id: "manage_recipients", label: "Получатели" },
  { id: "view_stats", label: "Статистика" },
  { id: "manage_integrations", label: "Интеграции" },
];

const emptySite = {
  site_name: "",
  full_name: "",
  tagline: "",
  description: "",
  region: "",
  departure_city: "",
  parent_org_name: "",
  parent_org_url: "",
  contact_phone: "",
  contact_phone_display: "",
  contact_email: "",
  mail_forward_to: "",
};

export default async function ManagementSettingsPage() {
  let unavailable = false;
  let channels: Array<{ id: string; configured: boolean; label: string }> = [];
  let events: Array<{
    kind: string;
    title: string;
    recipients: Array<{
      channel: string;
      address: string;
      event: string;
      ready: boolean;
      status: string;
      label: string;
    }>;
  }> = [];
  let site = emptySite;
  let roles: Array<{ id: string; name: string; permissions: string[] }> = [];
  let canManageRoles = false;

  try {
    const [session, notifications, siteSettings] = await Promise.all([
      getManagementSession(),
      getManagementNotificationSettings(),
      getManagementSiteSettings(),
    ]);
    canManageRoles = Boolean(session.full_admin || session.permissions?.includes("manage_roles"));
    channels = notifications.channels ?? [];
    events = notifications.events ?? [];
    site = { ...emptySite, ...siteSettings };
    if (canManageRoles) {
      roles = (await listManagementRoles()).map((role) => ({
        id: role.id,
        name: role.name,
        permissions: role.permissions ?? [],
      }));
    }
  } catch {
    unavailable = true;
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-stone-900">Настройки</h2>
        <p className="mt-1 max-w-2xl text-sm leading-6 text-stone-600">
          Идентичность сайта, получатели уведомлений (канал + адрес) и роли админки для полного
          администратора.
        </p>
      </div>
      {unavailable ? (
        <p className="rounded-2xl border border-stone-200 bg-white p-5 text-sm text-stone-600">
          Настройки сейчас недоступны. Проверьте вход и права роли.
        </p>
      ) : (
        <ManagementSettingsForms
          channels={channels}
          events={events}
          site={site}
          roles={roles}
          canManageRoles={canManageRoles}
          permissionOptions={permissionOptions}
        />
      )}
    </div>
  );
}
