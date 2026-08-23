"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { PassengerForm } from "@/components/account/passenger-form";
import type { Passenger, User } from "@/lib/api/auth";

type PassengerListProps = {
  user: User;
  passengers: Passenger[];
};

export function PassengerList({ user, passengers }: PassengerListProps) {
  const router = useRouter();
  const [editingID, setEditingID] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [deletingID, setDeletingID] = useState<string | null>(null);

  async function onDelete(id: string) {
    setDeletingID(id);
    setError(null);
    const response = await fetch(`/api/account/passengers/${id}`, { method: "DELETE" });
    setDeletingID(null);
    if (!response.ok) {
      const body = await response.json().catch(() => null);
      setError(body?.error ?? "Не удалось удалить пассажира");
      return;
    }
    if (editingID === id) {
      setEditingID(null);
    }
    router.refresh();
  }

  return (
    <div className="space-y-4">
      {error ? <p className="text-sm text-red-700">{error}</p> : null}
      {passengers.length === 0 ? (
        <p className="text-sm text-stone-500">Пока нет сохранённых пассажиров.</p>
      ) : (
        <ul className="space-y-3">
          {passengers.map((passenger) => (
            <li key={passenger.id} className="rounded-2xl border border-stone-200 bg-white p-5">
              {editingID === passenger.id ? (
                <PassengerForm
                  key={passenger.id}
                  user={user}
                  passenger={passenger}
                  onCancel={() => setEditingID(null)}
                />
              ) : (
                <div className="space-y-2">
                  <p className="font-medium text-stone-900">{passenger.name}</p>
                  <p className="text-sm text-stone-600">{passenger.phone}</p>
                  <p className="text-sm text-stone-600">Дата рождения: {passenger.birth_date}</p>
                  <p className="text-sm text-stone-600">Паспорт: {passenger.passport}</p>
                  <div className="flex flex-wrap gap-2 pt-2">
                    <button type="button" className="btn-secondary text-sm" onClick={() => setEditingID(passenger.id)}>
                      Изменить
                    </button>
                    <button
                      type="button"
                      className="btn-secondary text-sm"
                      disabled={deletingID === passenger.id}
                      onClick={() => onDelete(passenger.id)}
                    >
                      {deletingID === passenger.id ? "Удаляем…" : "Удалить"}
                    </button>
                  </div>
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
