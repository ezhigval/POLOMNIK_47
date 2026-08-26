"use client";

import { useState } from "react";
import { updateBookingPaymentStatusAction } from "@/app/management/actions";
import { formatManagementPaymentStatus } from "@/lib/format";

const PAYMENT_STATUSES = ["UNPAID", "AWAITING_PAYMENT", "PAID", "NOT_REQUIRED"] as const;

type BookingPaymentStatusFormProps = {
  bookingId: string;
  currentStatus: string;
};

export function BookingPaymentStatusForm({ bookingId, currentStatus }: BookingPaymentStatusFormProps) {
  const [status, setStatus] = useState(currentStatus);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);

  async function onSubmit() {
    setLoading(true);
    setError(null);
    setSaved(false);
    try {
      await updateBookingPaymentStatusAction(bookingId, status);
      setSaved(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось обновить статус оплаты");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
        <label className="sr-only" htmlFor={`payment-status-${bookingId}`}>
          Статус оплаты
        </label>
        <select
          id={`payment-status-${bookingId}`}
          value={status}
          onChange={(event) => {
            setStatus(event.target.value);
            setSaved(false);
          }}
          className="input-field py-2 text-sm"
        >
          {PAYMENT_STATUSES.map((value) => (
            <option key={value} value={value}>
              {formatManagementPaymentStatus(value)}
            </option>
          ))}
        </select>
        <button
          type="button"
          onClick={onSubmit}
          disabled={loading || status === currentStatus}
          className="btn-secondary px-4 py-2"
        >
          {loading ? "..." : "Сохранить"}
        </button>
      </div>
      {error ? <span className="text-sm text-red-700">{error}</span> : null}
      {saved ? <span className="text-sm text-emerald-700">Сохранено</span> : null}
    </div>
  );
}
