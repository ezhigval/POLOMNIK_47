"use client";

import { useState } from "react";
import { updateBookingStatusAction } from "@/app/management/actions";
import { formatManagementBookingStatus } from "@/lib/format";

const STATUSES = ["NEW", "CONTACTED", "CONFIRMED", "COMPLETED", "CANCELLED"] as const;

type BookingStatusFormProps = {
  bookingId: string;
  currentStatus: string;
};

export function BookingStatusForm({ bookingId, currentStatus }: BookingStatusFormProps) {
  const [status, setStatus] = useState(currentStatus);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);

  async function onSubmit() {
    setLoading(true);
    setError(null);
    setSaved(false);
    try {
      await updateBookingStatusAction(bookingId, status);
      setSaved(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось обновить статус");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
        <label className="sr-only" htmlFor={`status-${bookingId}`}>
          Статус заявки
        </label>
        <select
          id={`status-${bookingId}`}
          value={status}
          onChange={(event) => {
            setStatus(event.target.value);
            setSaved(false);
          }}
          className="input-field py-2 text-sm"
        >
          {STATUSES.map((value) => (
            <option key={value} value={value}>
              {formatManagementBookingStatus(value)}
            </option>
          ))}
        </select>
        <button
          type="button"
          onClick={onSubmit}
          disabled={loading || status === currentStatus}
          className="btn-primary px-4 py-2"
        >
          {loading ? "..." : "Сохранить"}
        </button>
      </div>
      {error ? <span className="text-sm text-red-700">{error}</span> : null}
      {saved ? <span className="text-sm text-emerald-700">Сохранено</span> : null}
    </div>
  );
}
