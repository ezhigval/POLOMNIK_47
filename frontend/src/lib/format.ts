const dateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "numeric",
  month: "long",
  year: "numeric",
});

const shortDateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "numeric",
  month: "short",
});

export function formatDateRange(start: string, end: string): string {
  const startDate = new Date(start);
  const endDate = new Date(end);
  if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) {
    return `${start} — ${end}`;
  }
  if (startDate.getFullYear() === endDate.getFullYear() && startDate.getMonth() === endDate.getMonth()) {
    return `${startDate.getDate()}–${dateFormatter.format(endDate)}`;
  }
  return `${shortDateFormatter.format(startDate)} — ${dateFormatter.format(endDate)}`;
}

export function formatPrice(amount: number, currency: string): string {
  const symbol = currency === "RUB" ? "₽" : currency;
  return `${amount.toLocaleString("ru-RU")} ${symbol}`;
}

export function formatTourDuration(start: string, end: string): string {
  const startDate = new Date(start);
  const endDate = new Date(end);
  if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) {
    return "";
  }
  const days = Math.round((endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24)) + 1;
  if (days <= 0) return "";
  if (days === 1) return "1 день";
  if (days >= 2 && days <= 4) return `${days} дня`;
  return `${days} дней`;
}

export function formatDateTime(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function formatBookingStatus(status: string): string {
  const labels: Record<string, string> = {
    NEW: "Принята",
    CONTACTED: "Менеджер свяжется с вами",
    CONFIRMED: "Подтверждена",
    COMPLETED: "Завершена",
    CANCELLED: "Отменена",
  };
  return labels[status] ?? status;
}

export function formatManagementBookingStatus(status: string): string {
  const labels: Record<string, string> = {
    NEW: "Новая",
    CONTACTED: "На связи",
    CONFIRMED: "Подтверждена",
    COMPLETED: "Завершена",
    CANCELLED: "Отменена",
  };
  return labels[status] ?? status;
}

export const integrationSyncLabels: Record<string, string> = {
  not_configured: "Не настроено",
  pending: "В очереди",
  synced: "Синхронизировано",
  failed: "Ошибка",
};

export const outboxStatusLabels: Record<string, string> = {
  pending: "Ожидает",
  processed: "Выполнено",
  failed: "Ошибка",
};

export type SlotsAvailability = "available" | "low" | "sold_out";

export function getSlotsAvailability(slotsLeft: number): SlotsAvailability {
  if (slotsLeft <= 0) return "sold_out";
  if (slotsLeft <= 3) return "low";
  return "available";
}

export function slotsLabel(slotsLeft: number): string {
  if (slotsLeft <= 0) return "Мест нет";
  if (slotsLeft === 1) return "Осталось 1 место";
  if (slotsLeft <= 4) return `Осталось ${slotsLeft} места`;
  return `Осталось ${slotsLeft} мест`;
}
