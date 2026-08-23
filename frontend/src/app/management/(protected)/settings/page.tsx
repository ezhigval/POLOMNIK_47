import { ManagementSettingsForms } from "@/components/management/management-settings-forms";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import {
  getManagementNotificationSettings,
  getManagementSession,
  getManagementSiteSettings,
  listManagementRoles,
  listManagementRoleTemplates,
} from "@/lib/api/management";
import { PERM, sessionHasPermission } from "@/lib/management-access";

const permissionOptions = [
  { id: "manage_tours", label: "Туры" },
  { id: "manage_bookings", label: "Заявки" },
  { id: "manage_support", label: "Поддержка" },
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
  const session = await getManagementSession().catch(() => null);
  if (!session) {
    return <ManagementNoAccess />;
  }

  const canManageSite = sessionHasPermission(session, PERM.settingsSite);
  const canManageRecipients = sessionHasPermission(session, PERM.recipients);
  const canManageRoles = sessionHasPermission(session, PERM.roles);
  if (!canManageSite && !canManageRecipients && !canManageRoles) {
    return <ManagementNoAccess />;
  }

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
  let roleTemplates: Array<{ id: string; label: string; role_name: string; permissions: string[] }> = [];

  if (canManageRecipients) {
    try {
      const notifications = await getManagementNotificationSettings();
      channels = notifications.channels ?? [];
      events = notifications.events ?? [];
    } catch {
      channels = [];
      events = [];
    }
  }
  if (canManageSite) {
    try {
      site = { ...emptySite, ...(await getManagementSiteSettings()) };
    } catch {
      site = emptySite;
    }
  }
  if (canManageRoles) {
    try {
      const [listed, templates] = await Promise.all([
        listManagementRoles(),
        listManagementRoleTemplates(),
      ]);
      roles = listed.map((role) => ({
        id: role.id,
        name: role.name,
        permissions: role.permissions ?? [],
      }));
      roleTemplates = templates;
    } catch {
      roles = [];
      roleTemplates = [];
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-stone-900">Настройки</h2>
        <p className="mt-1 max-w-2xl text-sm leading-6 text-stone-600">
          Идентичность сайта, получатели уведомлений (канал + адрес) и роли админки. Разделы зависят
          от прав текущей роли. Менять роли может только полный админ (`ADMIN_TOKEN`).
        </p>
      </div>
      <ManagementSettingsForms
        channels={channels}
        events={events}
        site={site}
        roles={roles}
        roleTemplates={roleTemplates}
        canManageRoles={canManageRoles}
        canManageSite={canManageSite}
        canManageRecipients={canManageRecipients}
        permissionOptions={permissionOptions}
      />
    </div>
  );
}
