const dateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "numeric",
  month: "long",
  year: "numeric",
});

const shortDateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "numeric",
  month: "short",
});

export function formatDateRange(start?: string | null, end?: string | null): string {
  if (!start || !end) {
    return "";
  }
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

export function formatPrice(amount: number | null | undefined, currency: string): string {
  if (amount == null) {
    return "";
  }
  const symbol = currency === "RUB" ? "₽" : currency;
  return `${amount.toLocaleString("ru-RU")} ${symbol}`;
}

export function formatTourDuration(start?: string | null, end?: string | null): string {
  if (!start || !end) {
    return "";
  }
  const startDate = new Date(start);
  const endDate = new Date(end);
  if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) {
    return "";
  }
  const days = Math.round((endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24)) + 1;
  if (days <= 0) return "";
  return `${days} ${pluralRu(days, "день", "дня", "дней")}`;
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
    CONTACTED: "Ожидает звонка",
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

export function formatPaymentStatus(status: string): string {
  const labels: Record<string, string> = {
    UNPAID: "Не оплачено",
    AWAITING_PAYMENT: "Ожидает оплаты",
    PAID: "Оплачено",
    NOT_REQUIRED: "Оплата не требуется",
  };
  return labels[status] ?? status;
}

export function formatManagementPaymentStatus(status: string): string {
  return formatPaymentStatus(status);
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
  return `Осталось ${slotsLeft} ${pluralRu(slotsLeft, "место", "места", "мест")}`;
}

/** Russian plural: 1 / 2–4 / 5–20, 21 / 22–24 / … */
export function pluralRu(n: number, one: string, few: string, many: string): string {
  const abs = Math.abs(n) % 100;
  const last = abs % 10;
  if (abs > 10 && abs < 20) return many;
  if (last === 1) return one;
  if (last >= 2 && last <= 4) return few;
  return many;
}

export function formatReviewCount(n: number): string {
  return `${n} ${pluralRu(n, "отзыв", "отзыва", "отзывов")}`;
}
