import { ManagementPanel } from "@/components/management/management-panel";
import { StatusBadge } from "@/components/management/status-badge";
import { getManagementMetricsDigest, getManagementWatchdog } from "@/lib/api/management";
import { formatDateTime } from "@/lib/format";
import { ManagementNoAccess } from "@/components/management/management-no-access";
import { canAccessManagementPage } from "@/lib/management-page-access";
import { PERM } from "@/lib/management-access";

const bookingStatuses = ["NEW", "CONTACTED", "CONFIRMED", "COMPLETED", "CANCELLED"] as const;

export default async function ManagementAIPage() {
  if (!(await canAccessManagementPage([PERM.stats]))) {
    return <ManagementNoAccess />;
  }

  const [digest, watchdog] = await Promise.all([
    getManagementMetricsDigest().catch(() => null),
    getManagementWatchdog().catch(() => null),
  ]);

  return (
    <div className="space-y-6">
      <ManagementPanel
        title="Дайджест метрик"
        description="Только заявки, опубликованные туры, открытые диалоги и outbox. Визиты не выдумываются."
      >
        {!digest ? (
          <p className="px-5 pb-5 text-sm text-stone-500">Не удалось загрузить дайджест.</p>
        ) : (
          <div className="space-y-4 px-5 pb-5">
            <p className="text-sm text-stone-600">
              AIPort: {digest.configured ? "настроен" : "noop — текст модели не строится, цифры из БД."}
            </p>
            <dl className="grid gap-3 sm:grid-cols-2">
              <div>
                <dt className="text-xs uppercase text-stone-500">Опубликованные туры</dt>
                <dd className="text-lg font-semibold">{digest.active_tours}</dd>
              </div>
              <div>
                <dt className="text-xs uppercase text-stone-500">Открытые диалоги</dt>
                <dd className="text-lg font-semibold">{digest.open_support_threads}</dd>
              </div>
              <div>
                <dt className="text-xs uppercase text-stone-500">Outbox pending / failed</dt>
                <dd className="text-lg font-semibold">
                  {digest.outbox_pending} / {digest.outbox_failed}
                </dd>
              </div>
            </dl>
            <div className="grid gap-2 sm:grid-cols-5">
              {bookingStatuses.map((status) => (
                <div key={status} className="rounded-lg bg-stone-50 px-3 py-2">
                  <p className="text-xs text-stone-500">{status}</p>
                  <p className="text-base font-semibold">{digest.bookings_by_status[status] ?? 0}</p>
                </div>
              ))}
            </div>
            {digest.summary ? (
              <p className="whitespace-pre-wrap text-sm leading-6 text-stone-700">{digest.summary}</p>
            ) : null}
          </div>
        )}
      </ManagementPanel>

      <ManagementPanel
        title="Watchdog"
        description="Отчёт: health, диск, outbox, 5xx с момента старта API, просроченный ночной бэкап (>26 ч). Рестарт прода не выполняется."
      >
        {!watchdog ? (
          <p className="px-5 pb-5 text-sm text-stone-500">Не удалось загрузить отчёт.</p>
        ) : (
          <div className="space-y-3 px-5 pb-5 text-sm">
            <div className="flex flex-wrap gap-2">
              <StatusBadge variant={watchdog.database === "ok" ? "success" : "warning"}>
                БД: {watchdog.database}
              </StatusBadge>
              <StatusBadge variant={watchdog.backup_overdue ? "warning" : "success"}>
                бэкап {watchdog.backup_overdue ? "просрочен" : "свежий"}
              </StatusBadge>
              <StatusBadge variant="neutral">рестарт: нет</StatusBadge>
            </div>
            <p className="text-stone-600">Снято: {formatDateTime(watchdog.at)}</p>
            <p className="text-stone-600">
              Диск: {watchdog.disk_percent}% ({watchdog.disk_used_bytes} / {watchdog.disk_total_bytes} байт)
            </p>
            <p className="text-stone-600">
              Outbox pending {watchdog.outbox_pending}, failed {watchdog.outbox_failed}; 5xx с старта API: {watchdog.status_5xx}
            </p>
            {watchdog.backup_at ? <p className="text-stone-600">Последний бэкап: {formatDateTime(watchdog.backup_at)}</p> : null}
            {watchdog.summary ? <p className="whitespace-pre-wrap text-stone-700">{watchdog.summary}</p> : null}
          </div>
        )}
      </ManagementPanel>
    </div>
  );
}
