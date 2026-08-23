"use client";

import { FormEvent, useMemo, useState } from "react";
import {
  updateNotificationSettingsAction,
  updateSiteSettingsAction,
  createAdminRoleAction,
  updateAdminRoleAction,
  deleteAdminRoleAction,
  assignAdminRoleUserAction,
} from "@/app/management/actions";
import { FormError } from "@/components/form-error";

export type SettingsChannel = { id: string; configured: boolean; label: string };
export type SettingsRecipient = {
  channel: string;
  address: string;
  event: string;
  ready: boolean;
  status: string;
  label: string;
};
export type SettingsEvent = {
  kind: string;
  title: string;
  recipients: SettingsRecipient[];
};
export type SiteSettingsFormData = {
  site_name: string;
  full_name: string;
  tagline: string;
  description: string;
  region: string;
  departure_city: string;
  parent_org_name: string;
  parent_org_url: string;
  contact_phone: string;
  contact_phone_display: string;
  contact_email: string;
  mail_forward_to: string;
};
export type AdminRoleFormData = {
  id: string;
  name: string;
  permissions: string[];
};
export type PermissionOption = { id: string; label: string };

const EVENT_KINDS = [
  { kind: "booking_created", title: "Новые заявки" },
  { kind: "booking_status_changed", title: "Смена статуса заявки" },
  { kind: "support_message", title: "Сообщения в поддержку" },
] as const;

type RecipientDraft = { channel: string; address: string };

type Props = {
  channels: SettingsChannel[];
  events: SettingsEvent[];
  site: SiteSettingsFormData;
  roles: AdminRoleFormData[];
  canManageRoles: boolean;
  canManageSite: boolean;
  canManageRecipients: boolean;
  permissionOptions: PermissionOption[];
};

function draftsFromEvents(events: SettingsEvent[]): Record<string, RecipientDraft[]> {
  const out: Record<string, RecipientDraft[]> = {};
  for (const meta of EVENT_KINDS) {
    const found = events.find((e) => e.kind === meta.kind);
    out[meta.kind] = (found?.recipients ?? []).map((r) => ({
      channel: r.channel || "telegram",
      address: r.address,
    }));
  }
  return out;
}

export function ManagementSettingsForms({
  channels,
  events,
  site,
  roles,
  canManageRoles,
  canManageSite,
  canManageRecipients,
  permissionOptions,
}: Props) {
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [drafts, setDrafts] = useState(() => draftsFromEvents(events));
  const [roleList, setRoleList] = useState(roles);

  const channelOptions = useMemo(
    () =>
      channels.length > 0
        ? channels
        : [
            { id: "telegram", configured: true, label: "Telegram" },
            { id: "max", configured: false, label: "Max" },
          ],
    [channels],
  );

  async function saveRecipients(event: FormEvent) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setSaved(null);
    try {
      const payload: Record<string, RecipientDraft[]> = {};
      for (const meta of EVENT_KINDS) {
        payload[meta.kind] = (drafts[meta.kind] ?? []).filter((r) => r.address.trim());
      }
      await updateNotificationSettingsAction(payload);
      setSaved("Получатели сохранены.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить получателей");
    } finally {
      setLoading(false);
    }
  }

  async function saveSite(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setSaved(null);
    const form = new FormData(event.currentTarget);
    try {
      await updateSiteSettingsAction({
        site_name: String(form.get("site_name") ?? ""),
        full_name: String(form.get("full_name") ?? ""),
        tagline: String(form.get("tagline") ?? ""),
        description: String(form.get("description") ?? ""),
        region: String(form.get("region") ?? ""),
        departure_city: String(form.get("departure_city") ?? ""),
        parent_org_name: String(form.get("parent_org_name") ?? ""),
        parent_org_url: String(form.get("parent_org_url") ?? ""),
        contact_phone: String(form.get("contact_phone") ?? ""),
        contact_phone_display: String(form.get("contact_phone_display") ?? ""),
        contact_email: String(form.get("contact_email") ?? ""),
        mail_forward_to: String(form.get("mail_forward_to") ?? ""),
      });
      setSaved("Настройки сайта сохранены.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось сохранить настройки сайта");
    } finally {
      setLoading(false);
    }
  }

  function updateDraft(kind: string, index: number, patch: Partial<RecipientDraft>) {
    setDrafts((prev) => {
      const list = [...(prev[kind] ?? [])];
      list[index] = { ...list[index], ...patch };
      return { ...prev, [kind]: list };
    });
  }

  function addDraft(kind: string) {
    setDrafts((prev) => ({
      ...prev,
      [kind]: [...(prev[kind] ?? []), { channel: "telegram", address: "" }],
    }));
  }

  function removeDraft(kind: string, index: number) {
    setDrafts((prev) => ({
      ...prev,
      [kind]: (prev[kind] ?? []).filter((_, i) => i !== index),
    }));
  }

  async function onCreateRole(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setSaved(null);
    const form = new FormData(event.currentTarget);
    const permissions = permissionOptions
      .map((p) => p.id)
      .filter((id) => form.get(`perm_${id}`) === "on");
    try {
      const created = await createAdminRoleAction({
        name: String(form.get("name") ?? ""),
        password: String(form.get("password") ?? ""),
        permissions,
      });
      setRoleList((prev) => [...prev, created]);
      event.currentTarget.reset();
      setSaved("Роль создана.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось создать роль");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-10">
      {canManageSite ? (
      <form onSubmit={saveSite} className="space-y-4 rounded-2xl border border-stone-200 bg-white p-5">
        <div>
          <h2 className="text-lg font-semibold text-stone-900">Идентичность сайта</h2>
          <p className="mt-1 text-sm text-stone-600">
            Публичные несекретные поля. Секреты (токены, ключи) только в env на сервере.
          </p>
        </div>
        <div className="grid gap-4 sm:grid-cols-2">
          {(
            [
              ["site_name", "Короткое имя", site.site_name],
              ["full_name", "Полное название", site.full_name],
              ["tagline", "Слоган", site.tagline],
              ["region", "Регион", site.region],
              ["departure_city", "Город выезда", site.departure_city],
              ["parent_org_name", "Родительская орг.", site.parent_org_name],
              ["parent_org_url", "URL орг.", site.parent_org_url],
              ["contact_phone", "Телефон (E.164)", site.contact_phone],
              ["contact_phone_display", "Телефон для показа", site.contact_phone_display],
              ["contact_email", "Email", site.contact_email],
            ] as const
          ).map(([name, label, value]) => (
            <label key={name} className="block text-sm">
              <span className="mb-1 block font-medium">{label}</span>
              <input name={name} defaultValue={value} className="input-field" />
            </label>
          ))}
        </div>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Описание</span>
          <textarea name="description" rows={3} defaultValue={site.description} className="input-field" />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block font-medium">Пересылка почты (через запятую)</span>
          <textarea
            name="mail_forward_to"
            rows={2}
            defaultValue={site.mail_forward_to}
            className="input-field"
            placeholder="smailikin70@yandex.ru"
          />
        </label>
        <button type="submit" className="btn-primary" disabled={loading}>
          Сохранить сайт
        </button>
      </form>
      ) : null}

      {canManageRecipients ? (
      <form onSubmit={saveRecipients} className="space-y-6 rounded-2xl border border-stone-200 bg-white p-5">
        <div>
          <h2 className="text-lg font-semibold text-stone-900">Получатели уведомлений</h2>
          <p className="mt-1 text-sm text-stone-600">
            Блоки канал + адрес. Telegram: username без @, человек должен написать боту /start. Max: телефон
            в международном формате (пока stub до появления ключей).
          </p>
          <ul className="mt-2 flex flex-wrap gap-3 text-xs text-stone-500">
            {channelOptions.map((ch) => (
              <li key={ch.id}>
                {ch.label}: {ch.configured ? "подключён" : "не настроен"}
              </li>
            ))}
          </ul>
        </div>

        {EVENT_KINDS.map((meta) => (
          <section key={meta.kind} className="space-y-3 border-t border-stone-100 pt-4">
            <div className="flex items-center justify-between gap-3">
              <h3 className="font-medium text-stone-900">{meta.title}</h3>
              <button type="button" className="text-sm text-brand-800" onClick={() => addDraft(meta.kind)}>
                + получатель
              </button>
            </div>
            {(drafts[meta.kind] ?? []).map((item, index) => {
              const status = events
                .find((e) => e.kind === meta.kind)
                ?.recipients.find((r) => r.address === item.address && r.channel === item.channel);
              return (
                <div key={`${meta.kind}-${index}`} className="grid gap-2 sm:grid-cols-[140px_1fr_auto]">
                  <select
                    className="input-field"
                    value={item.channel}
                    onChange={(e) => updateDraft(meta.kind, index, { channel: e.target.value })}
                  >
                    {channelOptions.map((ch) => (
                      <option key={ch.id} value={ch.id}>
                        {ch.label}
                      </option>
                    ))}
                  </select>
                  <input
                    className="input-field"
                    value={item.address}
                    onChange={(e) => updateDraft(meta.kind, index, { address: e.target.value })}
                    placeholder={item.channel === "max" ? "+79001234567" : "username"}
                  />
                  <button type="button" className="text-sm text-stone-500" onClick={() => removeDraft(meta.kind, index)}>
                    Удалить
                  </button>
                  {status ? (
                    <p className="sm:col-span-3 text-xs text-stone-500">
                      {status.status}
                      {status.ready ? " · готов" : ""}
                    </p>
                  ) : null}
                </div>
              );
            })}
          </section>
        ))}

        <button type="submit" className="btn-primary" disabled={loading}>
          Сохранить получателей
        </button>
      </form>
      ) : null}

      {canManageRoles ? (
        <section className="space-y-6 rounded-2xl border border-stone-200 bg-white p-5">
          <div>
            <h2 className="text-lg font-semibold text-stone-900">Роли админки</h2>
            <p className="mt-1 text-sm text-stone-600">
              Полный админ — только через ADMIN_TOKEN в env. Здесь создаются роли с паролем (хеш в БД) и
              правами; назначение — по UUID пользователя из кабинета.
            </p>
          </div>

          <form onSubmit={onCreateRole} className="space-y-3 border-b border-stone-100 pb-6">
            <div className="grid gap-3 sm:grid-cols-2">
              <label className="block text-sm">
                <span className="mb-1 block font-medium">Имя роли</span>
                <input name="name" required className="input-field" placeholder="manager" />
              </label>
              <label className="block text-sm">
                <span className="mb-1 block font-medium">Пароль (мин. 8)</span>
                <input name="password" type="password" required minLength={8} className="input-field" />
              </label>
            </div>
            <div className="flex flex-wrap gap-3">
              {permissionOptions.map((perm) => (
                <label key={perm.id} className="flex items-center gap-2 text-sm text-stone-700">
                  <input type="checkbox" name={`perm_${perm.id}`} />
                  {perm.label}
                </label>
              ))}
            </div>
            <button type="submit" className="btn-primary" disabled={loading}>
              Создать роль
            </button>
          </form>

          <ul className="space-y-4">
            {roleList.map((role) => (
              <li key={role.id} className="rounded-xl border border-stone-100 p-4">
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="font-medium text-stone-900">{role.name}</p>
                    <p className="mt-1 text-xs text-stone-500">{role.permissions.join(", ") || "без прав"}</p>
                  </div>
                  <button
                    type="button"
                    className="text-sm text-red-700"
                    onClick={async () => {
                      setLoading(true);
                      setError(null);
                      try {
                        await deleteAdminRoleAction(role.id);
                        setRoleList((prev) => prev.filter((r) => r.id !== role.id));
                        setSaved("Роль удалена.");
                      } catch (err) {
                        setError(err instanceof Error ? err.message : "Не удалось удалить роль");
                      } finally {
                        setLoading(false);
                      }
                    }}
                  >
                    Удалить
                  </button>
                </div>
                <form
                  className="mt-3 flex flex-wrap gap-2"
                  onSubmit={async (e) => {
                    e.preventDefault();
                    setLoading(true);
                    setError(null);
                    const form = new FormData(e.currentTarget);
                    try {
                      await assignAdminRoleUserAction(role.id, String(form.get("user_id") ?? ""));
                      setSaved("Пользователь привязан к роли.");
                      e.currentTarget.reset();
                    } catch (err) {
                      setError(err instanceof Error ? err.message : "Не удалось привязать UUID");
                    } finally {
                      setLoading(false);
                    }
                  }}
                >
                  <input
                    name="user_id"
                    required
                    className="input-field min-w-[280px] flex-1"
                    placeholder="UUID пользователя из /account"
                  />
                  <button type="submit" className="btn-primary" disabled={loading}>
                    Назначить
                  </button>
                </form>
                <form
                  className="mt-2 flex flex-wrap gap-2"
                  onSubmit={async (e) => {
                    e.preventDefault();
                    setLoading(true);
                    setError(null);
                    const form = new FormData(e.currentTarget);
                    try {
                      await updateAdminRoleAction(role.id, {
                        password: String(form.get("password") ?? ""),
                        permissions: role.permissions,
                      });
                      setSaved("Пароль роли обновлён.");
                      e.currentTarget.reset();
                    } catch (err) {
                      setError(err instanceof Error ? err.message : "Не удалось сменить пароль");
                    } finally {
                      setLoading(false);
                    }
                  }}
                >
                  <input
                    name="password"
                    type="password"
                    minLength={8}
                    required
                    className="input-field min-w-[200px] flex-1"
                    placeholder="Новый пароль роли"
                  />
                  <button type="submit" className="btn-primary" disabled={loading}>
                    Сменить пароль
                  </button>
                </form>
              </li>
            ))}
          </ul>
        </section>
      ) : null}

      {saved ? <p className="text-sm text-emerald-700">{saved}</p> : null}
      <FormError>{error}</FormError>
    </div>
  );
}
