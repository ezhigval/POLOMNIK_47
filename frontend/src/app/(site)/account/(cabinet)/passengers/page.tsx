import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { PassengerForm } from "@/components/account/passenger-form";
import { PassengerList } from "@/components/account/passenger-list";
import { fetchCurrentUser, fetchMyPassengers } from "@/lib/api/auth";
import { getAuthToken } from "@/lib/auth/session";

export const metadata: Metadata = {
  title: "Пассажиры",
};

export default async function AccountPassengersPage() {
  const token = await getAuthToken();
  if (!token) {
    redirect("/account/login?returnUrl=%2Faccount%2Fpassengers");
  }

  let user;
  try {
    user = await fetchCurrentUser(token);
  } catch {
    redirect("/account/login?returnUrl=%2Faccount%2Fpassengers");
  }

  const passengers = await fetchMyPassengers(token).catch(() => []);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900">Пассажиры</h1>
        <p className="mt-2 text-sm text-stone-600">
          Сохранённые спутники для заявок: ФИО, телефон, дата рождения и паспорт.
        </p>
      </div>

      <PassengerForm key={`${user.name}|${user.phone}|${passengers.length}`} user={user} />
      <PassengerList user={user} passengers={passengers} />
    </div>
  );
}
