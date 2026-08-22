import { TelegramSettingsForm } from "@/components/management/telegram-settings-form";
import { getManagementTelegramSettingsOrEmpty } from "@/lib/api/management";

export default async function ManagementSettingsPage() {
  const { settings, unavailable } = await getManagementTelegramSettingsOrEmpty();

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-stone-900">Настройки</h2>
        <p className="mt-1 max-w-2xl text-sm leading-6 text-stone-600">
          Получатели Telegram: отдельные списки для заявок и обращений в поддержку. Один и тот же
          бот, один токен.
        </p>
      </div>
      {unavailable || !settings ? (
        <p className="rounded-2xl border border-stone-200 bg-white p-5 text-sm text-stone-600">
          Настройки Telegram сейчас недоступны. Страница открыта; списки получателей появятся,
          когда API ответит.
        </p>
      ) : (
        <TelegramSettingsForm
          settings={{
            booking_usernames: settings.booking_usernames ?? "",
            support_usernames: settings.support_usernames ?? "",
            recipients: settings.recipients ?? [],
          }}
        />
      )}
    </div>
  );
}
