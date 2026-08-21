import {
  ManagementEmptyRow,
  ManagementPanel,
  ManagementTable,
  ManagementTableHead,
  ManagementTh,
} from "@/components/management/management-panel";
import { StatusBadge, syncStatusVariant } from "@/components/management/status-badge";
import {
  getManagementSystemInfo,
  listManagementIntegrationReferences,
  listManagementOutboxEvents,
} from "@/lib/api/management";
import { formatDateTime, integrationSyncLabels, outboxStatusLabels } from "@/lib/format";

function adapterBadge(adapter: string, configured: boolean) {
  if (adapter === "noop" || !adapter) {
    return <StatusBadge variant="neutral">noop (выкл.)</StatusBadge>;
  }
  if (configured) {
    return <StatusBadge variant="success">{adapter} · credentials OK</StatusBadge>;
  }
  return <StatusBadge variant="warning">{adapter} · ждёт credentials</StatusBadge>;
}

export default async function ManagementIntegrationsPage() {
  const [references, outboxEvents, systemInfo] = await Promise.all([
    listManagementIntegrationReferences(),
    listManagementOutboxEvents(),
    getManagementSystemInfo(),
  ]);

  const notificationEvents = outboxEvents.filter((event) =>
    event.event_type.startsWith("notification."),
  );
  const integrationEvents = outboxEvents.filter(
    (event) => !event.event_type.startsWith("notification."),
  );

  return (
    <div className="space-y-6">
      <ManagementPanel
        title="Адаптеры (код готов, подключение — финальный этап)"
        description="По умолчанию все внешние сервисы в режиме noop. Включение и отладка — после наполнения контентом и тестов."
      >
        <div className="grid gap-4 px-5 pb-5 sm:grid-cols-3">
          <div className="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p className="text-sm font-medium text-stone-900">Bitrix24 CRM</p>
            <div className="mt-2">{adapterBadge(systemInfo.crm_adapter, systemInfo.bitrix_configured)}</div>
          </div>
          <div className="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p className="text-sm font-medium text-stone-900">1С</p>
            <div className="mt-2">{adapterBadge(systemInfo.accounting_adapter, systemInfo.onec_configured)}</div>
          </div>
          <div className="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p className="text-sm font-medium text-stone-900">Telegram</p>
            <div className="mt-2">
              {adapterBadge(systemInfo.notification_adapter, systemInfo.telegram_configured)}
            </div>
          </div>
        </div>
      </ManagementPanel>

      <ManagementPanel
        title="Outbox здоровье"
        description="Сводка для мониторинга. Failed > 0 — повод смотреть worker и внешние адаптеры."
      >
        <div className="grid gap-4 px-5 pb-5 sm:grid-cols-3">
          <div className="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p className="text-sm text-stone-500">Pending</p>
            <p className="mt-1 text-2xl font-semibold text-stone-900">{systemInfo.outbox.pending}</p>
            {systemInfo.outbox.oldest_pending_at ? (
              <p className="mt-1 text-xs text-stone-500">
                Oldest: {formatDateTime(systemInfo.outbox.oldest_pending_at)}
              </p>
            ) : null}
          </div>
          <div className="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p className="text-sm text-stone-500">Failed</p>
            <p className="mt-1 text-2xl font-semibold text-stone-900">{systemInfo.outbox.failed}</p>
            {systemInfo.outbox.latest_failed_error ? (
              <p className="mt-1 line-clamp-2 text-xs text-red-700">{systemInfo.outbox.latest_failed_error}</p>
            ) : null}
          </div>
          <div className="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p className="text-sm text-stone-500">Processed</p>
            <p className="mt-1 text-2xl font-semibold text-stone-900">{systemInfo.outbox.processed}</p>
          </div>
        </div>
      </ManagementPanel>

      <ManagementPanel
        title="Синхронизация CRM / 1С"
        description="Записи появляются при CRM_ADAPTER=bitrix или ACCOUNTING_ADAPTER=onec. В noop таблицы пустые — это нормально."
      >
        <div className="px-5 pb-5">
          <p className="text-sm text-stone-600">
            Worker повторяет неудачные отправки в Bitrix24 и 1С через outbox.
          </p>
        </div>
      </ManagementPanel>

      <ManagementPanel title="Ссылки синхронизации">
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Сущность</ManagementTh>
            <ManagementTh>Система</ManagementTh>
            <ManagementTh>External ID</ManagementTh>
            <ManagementTh>Статус</ManagementTh>
            <ManagementTh>Обновлено</ManagementTh>
            <ManagementTh>Ошибка</ManagementTh>
          </ManagementTableHead>
          <tbody>
            {references.length === 0 ? (
              <ManagementEmptyRow colSpan={6}>Записей синхронизации пока нет.</ManagementEmptyRow>
            ) : (
              references.map((ref) => (
                <tr key={ref.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4">
                    <div className="font-medium">{ref.local_entity_type}</div>
                    <div className="font-mono text-xs text-stone-500">{ref.local_entity_id}</div>
                  </td>
                  <td className="px-4 py-4">
                    <div>{ref.external_system}</div>
                    <div className="text-stone-500">{ref.external_entity_type}</div>
                  </td>
                  <td className="px-4 py-4 font-mono text-xs">{ref.external_entity_id}</td>
                  <td className="px-4 py-4">
                    <StatusBadge variant={syncStatusVariant(ref.sync_status)}>
                      {integrationSyncLabels[ref.sync_status] ?? ref.sync_status}
                    </StatusBadge>
                  </td>
                  <td className="px-4 py-4 text-stone-600">{formatDateTime(ref.updated_at)}</td>
                  <td className="max-w-xs px-4 py-4 text-sm text-stone-600">{ref.last_error || "—"}</td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>

      <ManagementPanel title="Outbox — CRM / 1С">
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Тип</ManagementTh>
            <ManagementTh>Сущность</ManagementTh>
            <ManagementTh>Статус</ManagementTh>
            <ManagementTh>Попытки</ManagementTh>
            <ManagementTh>Создано</ManagementTh>
            <ManagementTh>Ошибка</ManagementTh>
          </ManagementTableHead>
          <tbody>
            {integrationEvents.length === 0 ? (
              <ManagementEmptyRow colSpan={6}>Событий в outbox пока нет.</ManagementEmptyRow>
            ) : (
              integrationEvents.map((event) => (
                <tr key={event.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4 font-mono text-xs">{event.event_type}</td>
                  <td className="px-4 py-4">
                    <div>{event.entity_type}</div>
                    <div className="font-mono text-xs text-stone-500">{event.entity_id}</div>
                  </td>
                  <td className="px-4 py-4">
                    <StatusBadge variant={syncStatusVariant(event.status)}>
                      {outboxStatusLabels[event.status] ?? event.status}
                    </StatusBadge>
                  </td>
                  <td className="px-4 py-4">{event.attempts}</td>
                  <td className="px-4 py-4 text-stone-600">{formatDateTime(event.created_at)}</td>
                  <td className="max-w-xs px-4 py-4 text-sm text-stone-600">{event.last_error || "—"}</td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>

      <ManagementPanel title="Outbox — Telegram уведомления">
        <ManagementTable>
          <ManagementTableHead>
            <ManagementTh>Тип</ManagementTh>
            <ManagementTh>Заявка</ManagementTh>
            <ManagementTh>Статус</ManagementTh>
            <ManagementTh>Попытки</ManagementTh>
            <ManagementTh>Создано</ManagementTh>
            <ManagementTh>Ошибка</ManagementTh>
          </ManagementTableHead>
          <tbody>
            {notificationEvents.length === 0 ? (
              <ManagementEmptyRow colSpan={6}>
                Событий Telegram пока нет (adapter noop или все доставлено).
              </ManagementEmptyRow>
            ) : (
              notificationEvents.map((event) => (
                <tr key={event.id} className="border-b border-stone-100 align-top last:border-0">
                  <td className="px-4 py-4 font-mono text-xs">{event.event_type}</td>
                  <td className="px-4 py-4 font-mono text-xs">{event.entity_id}</td>
                  <td className="px-4 py-4">
                    <StatusBadge variant={syncStatusVariant(event.status)}>
                      {outboxStatusLabels[event.status] ?? event.status}
                    </StatusBadge>
                  </td>
                  <td className="px-4 py-4">{event.attempts}</td>
                  <td className="px-4 py-4 text-stone-600">{formatDateTime(event.created_at)}</td>
                  <td className="max-w-xs px-4 py-4 text-sm text-stone-600">{event.last_error || "—"}</td>
                </tr>
              ))
            )}
          </tbody>
        </ManagementTable>
      </ManagementPanel>
    </div>
  );
}
